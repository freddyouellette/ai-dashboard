import { Button } from 'react-bootstrap';
import { useDispatch, useSelector } from 'react-redux';
import { goToBotEditPage, goToBotListPage, goToChatPage, selectSelectedBot } from '../store/page';
import { addOrUpdateBot, getBotModels, selectBotModels } from '../store/bots';
import RequiredStar from './RequiredStar';
import { persistChat } from '../store/chats';
import { useEffect } from 'react';
import moment from 'moment';

export default function CreateBotForm() {
	const botFormData = useSelector(selectSelectedBot) || {
		name: '',
		description: '',
		send_name: true,
		model: '',
		randomness: 1,
		personality: '',
		user_history: '',
	};
	const dispatch = useDispatch();
	
	const { botModels, botModelsLoading, botModelsError } = useSelector(selectBotModels);
	
	useEffect(() => {
		dispatch(getBotModels());
	}, [dispatch]);
	
	if (botModelsLoading) return <div>Loading...</div>;
	if (botModelsError) return <div>Error loading bot models...</div>;
	
	let handleChange = (event) => {
		let newBotFormData = {...botFormData};
		if (event.target.name === 'model') {
			let selectedOption = event.target.options[event.target.selectedIndex];
			let optGroupKey = selectedOption.parentNode.key;
			newBotFormData.ai_api_plugin_name = optGroupKey;
		}
		newBotFormData[event.target.name] = event.target.value;
		dispatch(goToBotEditPage(newBotFormData));
	}
	
	let handleSubmit = async (event) => {
		event.preventDefault();
		
		let createBotData = Object.assign({}, botFormData);
		createBotData.randomness = parseFloat(createBotData.randomness);
		let newBot = await dispatch(addOrUpdateBot(createBotData));
		if (createBotData.ID) {
			// updating existing bot
			dispatch(goToBotListPage())
		} else {
			let newChat = await dispatch(persistChat({
				name: 'New Chat',
				bot_id: newBot.ID,
			}));
			dispatch(goToChatPage(newChat));
		}
	}
	
	let botModelsList = Object.values(botModels);
	botModelsList.sort((a, b) => moment(b.created_at) - moment(a.created_at))
	
	let modelsByAuthor = {};
	for (let botModel of botModelsList) {
		modelsByAuthor[botModel.author_id] = modelsByAuthor[botModel.author_id] || [];
		modelsByAuthor[botModel.author_id].push(botModel);
	}
	let authorsAlphabetical = Object.keys(modelsByAuthor).sort();
	
	return (
		<div className="mx-3 mt-3">
			<h1 className="text-center">{botFormData?.ID ? 'Update Bot' : 'Create New Bot'}</h1>
			<form onSubmit={handleSubmit}>
				<div className="mb-3">
					<label htmlFor="name" className="form-label">Bot Name <RequiredStar/></label>
					<input onChange={handleChange} type="text" className="form-control" id="create-bot-form-name" name="name" placeholder="Enter bot name" required value={botFormData?.name ?? ''} />
				</div>
				<div className="mb-3 form-check">
					<input 
						onChange={e => handleChange({target: {name: 'send_name', value: e.target.checked}})} 
						className="form-check-input" 
						id="send_name" 
						name="send_name" 
						type="checkbox"
						checked={botFormData?.send_name}
					/>
					<label htmlFor="send_name" className="form-check-label text-left d-flex justify-content-start">Tell the bot its name</label>
				</div>
				<div className="mb-3">
					<label htmlFor="description" className="form-label">Description</label>
					<input onChange={handleChange} type="text" className="form-control" id="create-bot-form-description" name="description" placeholder="Enter bot description" value={botFormData?.description ?? ''} />
				</div>
				<div className="mb-3">
					<label htmlFor="model" className="form-label">Model <RequiredStar/></label>
					<select onChange={handleChange} id="create-bot-form-model" name="model" className="form-control" value={botFormData?.model ?? 'gpt-4'}>
						<option value="">Select AI Model</option>
						{authorsAlphabetical.map(author_id => {
							return <optgroup key={author_id} label={modelsByAuthor[author_id][0].author_name}>
								{modelsByAuthor[author_id].map(botModel => {
									return <option key={botModel.id} value={botModel.id}>{botModel.id}</option>
								})}
							</optgroup>
						})}
					</select>
				</div>
				<div className="mb-3">
					<label htmlFor="personality" className="form-label">Personality</label>
					<textarea onChange={handleChange} className="form-control" id="create-bot-form-personality" name="personality" rows="3" placeholder="Enter bot personality" value={botFormData?.personality ??  ''} ></textarea>
				</div>
				<div className="mb-3">
					<label htmlFor="user_history" className="form-label">User History</label>
					<textarea onChange={handleChange} className="form-control" id="create-bot-form-user_history" name="user_history" rows="3" placeholder="Enter user history" value={botFormData?.user_history ??  ''} ></textarea>
				</div>
				<div className="mb-3">
					<label htmlFor="correction_prompt" className="form-label">Message Correction Prompt</label>
					<textarea onChange={handleChange} className="form-control" id="create-bot-form-correction_prompt" name="correction_prompt" rows="3" placeholder="Enter Correction Prompt" value={botFormData?.correction_prompt ?? ''} ></textarea>
					<div className="help-text text-start">This prompt will be used to correct each message you send. Leave blank to turn off corrections.</div>
				</div>
				<div className="mb-3">
					<label htmlFor="randomness" className="form-label">Randomness <RequiredStar/></label>
					<input 
						onChange={handleChange} 
						className="form-control" 
						id="create-bot-form-randomness" 
						name="randomness" 
						placeholder="Randomness"
						type="number"
						min="0"
						max="1"
						step="any"
						required
						value={botFormData?.randomness ?? 1} />
					<small><em>Between 0 and 1. 0 = deterministic. 1 = very random.</em></small>
				</div>
				<Button className="btn-secondary me-3" onClick={() => dispatch(goToBotListPage())}>Cancel</Button>
				<button type="submit" className="btn btn-primary">{botFormData?.ID ? 'Update' : 'Create'} Bot</button>
			</form>
		</div>
	);
}