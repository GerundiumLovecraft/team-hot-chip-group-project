/*
Allow logged in users to submit a new coffee spot
- Redirect to login page if logged out
- Form fields: name, address, description, open from, open to, url image
- Form fields must not be empty
- Feature selection: fetch features from backend and display as checkboxes
- Clicking submit takes you to feed page and card will display

spotResult := JSONSpot{
		ID:          spotFound.ID,
		UserId:      spotFound.UserId,
		Name:        spotFound.Name,
		Address:     spotFound.Address,
		Description: spotFound.Description,
		OpenFrom:    spotFound.OpenFrom,
		OpenTo:      spotFound.OpenTo,
		Features:    jsonFeature,
	}
 */

import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom"
import { createSpot } from "../../services/spots.js"
import { getFeatures }  from "../../services/features.js";
import {isTokenValid} from "../../helpers/authentication.js";
import AnimatedButton from "../../components/animation/AnimatedButton";
import "./CreateSpotPage.css"

export function CreateSpotPage() {
    const [name, setName] = useState("");
    const [address, setAddress] = useState("");
    const [description, setDescription] = useState("");
    const [image, setImage] = useState("")
    const [openFrom, setOpenFrom] = useState("");
    const [openTo, setOpenTo] = useState("");
    const [listedFeatures, setListedFeatures] = useState([])
    const [addFeatures, setAddFeatures] = useState([]);
    const [errors, setErrors] = useState([]);
    const token = localStorage.getItem("token")
    const navigate = useNavigate();

    const featureDict = {
        "noise_level": ["Quiet", "Moderate", "Loud"],
        "price": ["£", "££", "£££"],
    }

    const featureNames = {
    "wifi": "WiFi",
    "toilets": "Toilets",
    "power_sockets": "Power Sockets",
    "open_late": "Open Late",
    }

    useEffect(() => {
        if (!isTokenValid()) {
            navigate("/login");
            return;
        }

        getFeatures()
        .then((data) => {
            setListedFeatures(data.features);
        })
        .catch((err) => {
            console.error(err);
            navigate("/");
        });

    }, [navigate])

    function handleChange(event) {
        const inputValue = event.target.value;
        const inputName = event.target.name

        switch(inputName) {
            case "name":
                setName(inputValue);
                break;
            case "address":
                setAddress(inputValue);
                break;
            case "description":
                setDescription(inputValue);
                break;
            case "image":
                setImage(inputValue);
                break;
            case "open-from":
                setOpenFrom(inputValue);
                break;
            case "open-to":
                setOpenTo(inputValue);
                break;
            case "features-no-value": {
                const changedFeat = {
                    "feat_id": event.target.dataset.featId,
                    "value": null,
                };

                if (addFeatures.some((feat) => feat.feat_id === changedFeat.feat_id)) {
                    setAddFeatures(addFeatures.filter((feat) => feat.feat_id !== changedFeat.feat_id));
                } else {
                    setAddFeatures([...addFeatures, changedFeat]);
                }

                break;
            }
            case "features-with-value": {
                const changedFeat = {
                    "feat_id": event.target.dataset.featId,
                    "value": event.target.dataset.featValue,
                };

                const filtered = addFeatures.filter((feat) => feat.feat_id !== changedFeat.feat_id);
                if (addFeatures.some((feat) => feat.feat_id === changedFeat.feat_id && feat.value === changedFeat.value)) {
                    setAddFeatures(filtered);
                } else {
                    setAddFeatures([...filtered, changedFeat]);
                }

                break;
            }
            default:
                break;
        }
    }

    function checkErrors(spotJson) {
        let newErrors = []

        if (spotJson.name.length === 0) {
            newErrors.push("Name field should not be empty");
        }
        if (spotJson.address.length === 0) {
            newErrors.push("Address field should not be empty");
        }
        if (spotJson.description.length === 0) {
            newErrors.push("Description field should not be empty");
        }
        if (spotJson.image.length === 0) {
            newErrors.push("Provide a URL for the image");
        } else {
            try {
                new URL(spotJson.image)
            } catch (e) {
                newErrors.push("Provide a valid URL");
            }
        }
        if (spotJson.open_from.length !== 5) {
            newErrors.push("Provide a correct open time");
        }
        if (spotJson.open_to.length !== 5) {
            newErrors.push("Provide a correct close time");
        }

        const hasNoise = addFeatures.some(f => {
            const feature = listedFeatures.find(lf => String(lf.feat_id) === f.feat_id);
            return feature?.feat_name === "noise_level";
        });
        const hasPrice = addFeatures.some(f => {
            const feature = listedFeatures.find(lf => String(lf.feat_id) === f.feat_id);
            return feature?.feat_name === "price";
        });
        const hasFeature = addFeatures.some(f => {
            const feature = listedFeatures.find(lf => String(lf.feat_id) === f.feat_id);
            return feature?.feat_name !== "noise_level" && feature?.feat_name !== "price" && feature?.feat_name !== "open_late";
        });

        if (!hasNoise) newErrors.push("Select a noise level");
        if (!hasPrice) newErrors.push("Select a price range");
        if (!hasFeature) newErrors.push("Select at least one feature");

        setErrors(newErrors);
        return newErrors.length > 0;
    }

    async function handleSubmit(event) {
        event.preventDefault();

        try {
            const openLateFeat = listedFeatures.find((f) => f.feat_name === "open_late");
            const finalFeatures = (openTo >= "18:00" || openTo < openFrom)
                ? [...addFeatures, { "feat_id": openLateFeat.feat_id, "value": null }]
                : addFeatures;

            const newSpotJson = {
                "name": name,
                "address": address,
                "description": description,
                "image": image,
                "open_from": openFrom,
                "open_to": openTo,
                "features": finalFeatures.map((f) => ({
                    "feat_id": Number(f.feat_id),
                    "value": f.value !== null ? parseInt(f.value, 10) : null
                })),
            };

            if (checkErrors(newSpotJson)) {
                return;
            }

            await createSpot(newSpotJson, token);

            setErrors([]);
            setName("");
            setAddress("");
            setDescription("");
            setImage("");
            setOpenFrom("");
            setOpenTo("");
            setAddFeatures([]);
            navigate("/");

        } catch(e) {
            setErrors(["Something went wrong, please try again!"]);
        }
    }

    return(
        <>
            <div className="create-spot-page">
                <div className="create-spot-header">
                    <h2>Submit a New Spot</h2>
                    <p>Share your favourite work spot with the community</p>
                </div>
                <form onSubmit={handleSubmit} className="create-spot-form">

                    {errors.length > 0 && (
                        <div className="form-errors">
                            {errors.map((e, i) => <p key={i}>{e}</p>)}
                        </div>
                    )}

                    <div className="form-field">
                        <label>Spot name <span className="required">*</span></label>
                        <input
                            name="name"
                            type="text"
                            placeholder="What's the name of the cafe?"
                            value={name}
                            onChange={(e) => handleChange(e)}
                        />
                    </div>
                    <div className="form-field">
                        <label>Address <span className="required">*</span></label>
                        <input
                            name="address"
                            type="text"
                            placeholder="e.g. 12 Baker Street, London, EC1A 1BB"
                            value={address}
                            onChange={(e) => handleChange(e)}
                        />
                    </div>
                    <div className="form-field full-width">
                        <label>Description <span className="required">*</span></label>
                        <textarea
                            name="description"
                            placeholder="Tell us what makes this spot great for working..."
                            value={description}
                            onChange={(e) => handleChange(e)}
                        />
                    </div>
                    <div className="form-field full-width">
                        <label>Image URL <span className="required">*</span></label>
                        <input
                            name="image"
                            placeholder="Paste a link to the image here"
                            value={image}
                            onChange={(e) => handleChange(e)}
                        />
                    </div>
                    <div className="time-row">
                        <div className="form-field">
                            <label>Open from <span className="required">*</span></label>
                            <input
                                name="open-from"
                                type="time"
                                value={openFrom}
                                onChange={(e) => handleChange(e)}
                            />
                        </div>
                        <div className="form-field">
                            <label>Open to <span className="required">*</span></label>
                            <input
                                name="open-to"
                                type="time"
                                value={openTo}
                                onChange={(e) => handleChange(e)}
                            />
                        </div>
                    </div>

                    <hr className="features-divider" />

                    {listedFeatures.map((feature) => {
                        if (feature.feat_name === "noise_level" || feature.feat_name === "price") {
                            return (
                                <div key={feature.feat_id} className="scale-feature">
                                    <label>{feature.feat_name === "noise_level" ? "Noise Level" : "Price"} <span className="required">*</span></label>
                                    <div className="scale-buttons">
                                        <button
                                            type="button"
                                            className={`scale-btn ${addFeatures.some(f => f.feat_id === String(feature.feat_id) && f.value === "1") ? "selected" : ""}`}
                                            name="features-with-value"
                                            data-feat-id={feature.feat_id}
                                            data-feat-value="1"
                                            onClick={(e) => handleChange(e)}
                                        >
                                            {featureDict[feature.feat_name][0]}
                                        </button>
                                        <button
                                            type="button"
                                            className={`scale-btn ${addFeatures.some(f => f.feat_id === String(feature.feat_id) && f.value === "2") ? "selected" : ""}`}
                                            name="features-with-value"
                                            data-feat-id={feature.feat_id}
                                            data-feat-value="2"
                                            onClick={(e) => handleChange(e)}
                                        >
                                            {featureDict[feature.feat_name][1]}
                                        </button>
                                        <button
                                            type="button"
                                            className={`scale-btn ${addFeatures.some(f => f.feat_id === String(feature.feat_id) && f.value === "3") ? "selected" : ""}`}
                                            name="features-with-value"
                                            data-feat-id={feature.feat_id}
                                            data-feat-value="3"
                                            onClick={(e) => handleChange(e)}
                                        >
                                            {featureDict[feature.feat_name][2]}
                                        </button>
                                    </div>
                                </div>
                            )
                        } else if (feature.feat_name !== "open_late") {
                            return null
                        }
                    })}

                    <div className="features-section">
                        <label>Features <span className="required">*</span></label>
                        <div className="features-grid">
                            {listedFeatures
                                .filter(f => f.feat_name !== "noise_level" && f.feat_name !== "price" && f.feat_name !== "open_late")
                                .map((feature) => (
                                    <button
                                        key={feature.feat_id}
                                        type="button"
                                        className={`feature-btn ${addFeatures.some(f => f.feat_id === String(feature.feat_id)) ? "selected" : ""}`}
                                        name="features-no-value"
                                        data-feat-id={feature.feat_id}
                                        value={feature.feat_id}
                                        onClick={(e) => handleChange(e)}
                                    >
                                        {featureNames[feature.feat_name] || feature.feat_name}
                                    </button>
                                ))
                            }
                        </div>
                    </div>

                    <AnimatedButton 
                        type="submit" className="submit-btn">Create spot!
                    </AnimatedButton>
                </form>
            </div>
        </>
    )
}