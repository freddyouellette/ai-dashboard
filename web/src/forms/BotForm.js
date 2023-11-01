import axios from 'axios';
import { Button } from 'react-bootstrap';
import { useDispatch, useSelector } from 'react-redux';
import { selectBotToUpdate, goToBotEdit } from '../store/page';
import { addOrUpdateBot } from '../store/bots';
import RequiredStar from './RequiredStar';

export default function CreateBotForm() {
	const botFormData = useSelector(selectBotToUpdate) || {
		name: '',
		description: '',
		model: 'gpt-4',
		randomness: 1,
		personality: '',
		user_history: '',
	};
	const dispatch = useDispatch();
	
	let handleChange = (event) => {
		dispatch(goToBotEdit({
			...botFormData,
			[event.target.name]: event.target.value
		}));
	}
	
	let handleSubmit = async (event) => {
		event.preventDefault();
		
		let createBotData = Object.assign({}, botFormData);
		createBotData.randomness = parseFloat(createBotData.randomness);
		axios[botFormData.ID ? "put" : "post"]("http://localhost:8080/api/bots", createBotData)
		.then(response => {
			console.log(response);
			dispatch(addOrUpdateBot(response.data));
		}).catch(error => {
			console.error(error);
		})
	}
	
	return (
		<div className="mx-3">
			<h1 className="text-center">{botFormData?.ID ? 'Update Bot' : 'Create New Bot'}</h1>
			<form onSubmit={handleSubmit}>
				<div className="mb-3">
					<label htmlFor="name" className="form-label">Bot Name <RequiredStar/></label>
					<input onChange={handleChange} type="text" className="form-control" id="create-bot-form-name" name="name" placeholder="Enter bot name" required value={botFormData?.name ?? ''} />
				</div>
				<div className="mb-3">
					<label htmlFor="description" className="form-label">Description</label>
					<input onChange={handleChange} type="text" className="form-control" id="create-bot-form-description" name="description" placeholder="Enter bot description" value={botFormData?.description ?? ''} />
				</div>
				<div className="mb-3">
					<label htmlFor="model" className="form-label">Model <RequiredStar/></label>
					<select onChange={handleChange} id="create-bot-form-model" name="model" className="form-control" value={botFormData?.model ?? 'gpt-4'}>
						<option value="">Select AI Model</option>
						<option value="gpt-4">GPT-4</option>
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
					<label htmlFor="randomness" className="form-label">Randomness <RequiredStar/></label>
					<input 
						onChange={handleChange} 
						className="form-control" 
						id="create-bot-form-randomness" 
						name="randomness" 
						placeholder="Randomness"
						type="number"
						min="0"
						max="2"
						step="any"
						required
						value={botFormData?.randomness ?? 1} />
					<small><em>Between 0 and 2. 0 = deterministic. 2 = very random</em></small>
				</div>
				<Button className="btn-secondary me-3">Cancel</Button>
				<button type="submit" className="btn btn-primary">{botFormData?.ID ? 'Update' : 'Create'} Bot</button>
			</form>
		</div>
	);
}