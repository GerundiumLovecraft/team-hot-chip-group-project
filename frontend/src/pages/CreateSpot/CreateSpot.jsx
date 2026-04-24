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

function CreateSpotPage() {
    const [name, setName] = useState("");
    const [address, setAddress] = useState("");
    const [description, setDescription] = useState("");
    const [openFrom, setOpenFrom] = useState("");
    const [openTo, setOpenTo] = useState("");
    const [listedFeatures, setListedFeatures] = useState([])
    const [addFeatures, setAddFeatures] = useState([]);
    const [errors, setErrors] = useState([]);
    const token = localStorage.getItem("token")
    const navigate = useNavigate();

    const featureDict = {
        "noise_level": ["quiet", "moderate", "loud"],
        "price": ["£", "££", "£££"],
    }

    useEffect(() => {
        if (!token) {
            navigate("/login");
            return;
        }

        getFeatures().
        then((data) => {
            setListedFeatures(data.features);
        }).
        catch((err) => {
            console.error(err);
            navigate("/login")
        })

    }, [navigate])

    function handleChange(event) {
        const inputValue = event.target.value;
        const inputName = event.target.name

        switch(inputName) {
            case "name":
                setName(inputValue)
                break;
            case "address":
                setAddress(inputValue)
                break;
            case "description":
                setDescription(inputValue)
                break;
            case "open-from":
                setOpenFrom(inputValue)
                break;
            case "open-to":
                setOpenTo(inputValue)
                break;
            case "features-no-value": {
                const changedFeat = {
                    "feat_name": event.target.dataset.featName,
                    "value": null,
                };

                if (addFeatures.some((feat) => feat.feat_name === changedFeat.feat_name)) {
                    setAddFeatures(addFeatures.filter((feat) => feat.feat_name != changedFeat.feat_name));
                } else {
                    setAddFeatures([...addFeatures, changedFeat]);
                }
                break;
            }
            case "features-with-value": {
                const changedFeat = {
                    "feat_name": event.target.dataset.featName,
                    "value": event.target.dataset.featValue,
                };
                const filtered = addFeatures.filter((feat) => feat.feat_name !== changedFeat.feat_name);
                if (addFeatures.some((feat) => feat.feat_name === changedFeat.feat_name && feat.value === changedFeat.value)) {
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

    async function handleSubmit(event) {
        event.preventDefault();

        try{
            const finalFeatures = (openTo >= "18:00" || openTo < openFrom)
                ? [...addFeatures, { "feat_name": "open_late", "value": null }]
                : addFeatures;

            const newSpotJson = {
                "name": name,
                "address": address,
                "description": description,
                "open_from": openFrom,
                "open_to": openTo,
                "features": finalFeatures,
            };

            await createSpot(newSpotJson, token)
            navigate("/browse_spots")
        } catch {
            setErrors([...errors, "Something went wrong, please try again!"])
        } finally {
            setName("")
            setAddress("")
            setDescription("")
            setOpenFrom("")
            setOpenTo("")
            setAddFeatures([])
        }
    }

    return(
        <>
            <div>
                <h2>Create a new Spot</h2>
                {(errors.length > 0) && <p className={errors}></p>}
                <form onSubmit={handleSubmit}>
                    <input
                        name="name"
                        type="text"
                        placeholder="How's your spot called?"
                        value={name}
                        onChange={(e) => handleChange(e)}
                        required
                    />
                    <input
                        name="address"
                        type="text"
                        placeholder="Where's your spot located"
                        value={address}
                        onChange={(e) => {handleChange(e)}}
                        required
                    />
                    <textarea
                        name="description"
                        placeholder="Write a short description about your spot..."
                        value={description}
                        onChange={(e) => {handleChange(e)}}
                        required
                    />
                    <input
                        name="open-from"
                        type="time"
                        value={openFrom}
                        onChange={(e) => {handleChange(e)}}
                        required
                    />
                    <input
                        name="open-to"
                        type="time"
                        value={openTo}
                        onChange={(e) => {handleChange(e)}}
                        required
                    />
                    {listedFeatures.map((feature) => {
                        if (feature.feature_name === "noise_level" || feature.feature_name === "price") {
                            return (
                                <div key={feature.feat_id}>
                                    <label>{feature.feature_name}</label>
                                    <button
                                        name="features-with-value"
                                        data-feat-name={feature.feature_name}
                                        data-feat-value={1}
                                        onClick={(e) => {handleChange(e)}}
                                    >
                                        {featureDict[feature.feature_name][0]}
                                    </button>
                                    <button
                                        name="features-with-value"
                                        data-feat-name={feature.feature_name}
                                        data-feat-value={2}
                                        onClick={(e) => {handleChange(e)}}
                                    >
                                        {featureDict[feature.feature_name][1]}
                                    </button>
                                    <button
                                        name="features-with-value"
                                        data-feat-name={feature.feature_name}
                                        data-feat-value={3}
                                        onClick={(e) => {handleChange(e)}}
                                    >
                                        {featureDict[feature.feature_name][2]}
                                    </button>
                                </div>
                            )
                        } else {
                            return (
                                <button
                                    key={feature.feat_id}
                                    name="features-no-value"
                                    data-feat-name={feature.feature_name}
                                    value={feature.feature_name}
                                    onClick={(e) => {handleChange(e)}}
                                />
                            )
                        }
                    })}
                </form>
            </div>
        </>
    )
}