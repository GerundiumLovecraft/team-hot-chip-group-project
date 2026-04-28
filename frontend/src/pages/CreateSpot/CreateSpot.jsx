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

    useEffect(() => {
        if (!isTokenValid()) {
            navigate("/login");
            return;
        }

        getFeatures().
        then((data) => {
            setListedFeatures(data.features);
        }).
        catch((err) => {
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

        let errorCount = 0;
        let newErrors = []

        if (spotJson.name.length === 0) {
            newErrors = [...newErrors, "Name field should not be empty"]
            errorCount++;
        }
        if (spotJson.address.length === 0) {
            newErrors = [...newErrors, "Address field should not be empty"]
            errorCount++;
        }
        if (spotJson.description.length === 0) {
            newErrors = [...newErrors, "Description field should not be empty"];
            errorCount++;
        }
        if (spotJson.image.length === 0) {
            newErrors = [...newErrors, "Provide a URL for the image"]
            errorCount++;
        } else {
            try {
                new URL(spotJson.image)
            } catch (e) {
                newErrors = [...newErrors, "Provide a valid URL"];
                errorCount++;
            }
        }
        if (spotJson.open_from.length !== 5) {
            newErrors = [...newErrors, "Provide a correct open time"];
            errorCount++;
        }
        if (spotJson.open_to.length !== 5) {
            newErrors = [...newErrors, "Provide a correct close time"];
            errorCount++;
        }
        if (spotJson.features.length === 0) {
            newErrors = [...newErrors, "Select at least one feature"];
            errorCount++;
        }

        if (errorCount > 0) {
            setErrors(newErrors)
            return true;
        } else {
            return false;
        }
    }

    async function handleSubmit(event) {
        event.preventDefault();

        try{
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
                    "value": f.value !== null ? parseInt(f.value, 10) : null})),
            };

            if (checkErrors(newSpotJson)) {
                throw new Error("Validator errors");
            }

            await createSpot(newSpotJson, token);
            setErrors([]);
            navigate("/");
        } catch(e) {
            if (e.message !== "Validator errors") {
                setErrors([...errors, "Something went wrong, please try again!"])
            }
        } finally {
            setName("");
            setAddress("");
            setDescription("");
            setImage("");
            setOpenFrom("");
            setOpenTo("");
            setAddFeatures([]);
        }
    }

    return(
        <>
            <div className="create-spot-page">
                <h2>Create a new Spot</h2>
                <div className="form-errors">
                    {errors.map((e, i) => <p key={i}>{e}</p>)}
                </div>
                <form onSubmit={handleSubmit}>
                    <div className="form-field">
                        <label>Spot name</label>
                        <input
                            name="name"
                            type="text"
                            placeholder="How's your spot called?"
                            value={name}
                            onChange={(e) => handleChange(e)}
                        />
                    </div>
                    <div className="form-field">
                        <label>Address</label>
                        <input
                            name="address"
                            type="text"
                            placeholder="Where's your spot located"
                            value={address}
                            onChange={(e) => {handleChange(e)}}
                        />
                    </div>
                    <div className="form-field">
                        <label>Description</label>
                        <textarea
                            name="description"
                            placeholder="Write a short description about your spot..."
                            value={description}
                            onChange={(e) => {handleChange(e)}}
                        />
                    </div>
                    <div className="form-field">
                        <label>Image URL</label>
                        <input
                            name="image"
                            placeholder="Paste a link to the image here"
                            value={image}
                            onChange={(e) => {handleChange(e)}}
                        />
                    </div>
                    <div className="time-row">
                        <div className="form-field">
                            <label>Open from</label>
                            <input
                                name="open-from"
                                type="time"
                                value={openFrom}
                                onChange={(e) => {handleChange(e)}}
                            />
                        </div>
                        <div className="form-field">
                            <label>Open to</label>
                            <input
                                name="open-to"
                                type="time"
                                value={openTo}
                                onChange={(e) => {handleChange(e)}}
                            />
                        </div>
                    </div>

                    {listedFeatures.map((feature) => {
                        if (feature.feat_name === "noise_level" || feature.feat_name === "price") {
                            return (
                                <div key={feature.feat_id}>
                                    <label>{feature.feat_name}</label>
                                    <button
                                        type="button"
                                        className={`scale-btn ${addFeatures.some(f => f.feat_id === String(feature.feat_id) && f.value === "1") ? "selected" : ""}`}
                                        name="features-with-value"
                                        data-feat-id={feature.feat_id}
                                        data-feat-value="1"
                                        onClick={(e) => {handleChange(e)}}
                                    >
                                        {featureDict[feature.feat_name][0]}
                                    </button>
                                    <button
                                        type="button"
                                        className={`scale-btn ${addFeatures.some(f => f.feat_id === String(feature.feat_id) && f.value === "2") ? "selected" : ""}`}
                                        name="features-with-value"
                                        data-feat-id={feature.feat_id}
                                        data-feat-value="2"
                                        onClick={(e) => {handleChange(e)}}
                                    >
                                        {featureDict[feature.feat_name][1]}
                                    </button>
                                    <button
                                        type="button"
                                        className={`scale-btn ${addFeatures.some(f => f.feat_id === String(feature.feat_id) && f.value === "3") ? "selected" : ""}`}
                                        name="features-with-value"
                                        data-feat-id={feature.feat_id}
                                        data-feat-value="3"
                                        onClick={(e) => {handleChange(e)}}
                                    >
                                        {featureDict[feature.feat_name][2]}
                                    </button>
                                </div>
                            )
                        } else if (feature.feat_name != "open_late") {
                            return (
                                <button
                                    key={feature.feat_id}
                                    type="button"
                                    className={`feature-btn ${addFeatures.some(f => f.feat_id === String(feature.feat_id)) ? "selected" : ""}`}
                                    name="features-no-value"
                                    data-feat-id={feature.feat_id}
                                    value={feature.feat_id}
                                    onClick={(e) => {handleChange(e)}}
                                >
                                    {feature.feat_name}
                                </ button>
                            )
                        }
                    })}
                    <button type="submit" className="submit-btn"> Create a spot!</button>
                </form>
            </div>
        </>
    )
}