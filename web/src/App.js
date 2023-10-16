import { Button, Col, Container, ListGroup, ListGroupItem, Row } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faPlus } from '@fortawesome/free-solid-svg-icons'
import './App.css';
import { useState } from 'react';

function App() {
	return (
		<Row className="App h-100">
			<Col className="sidebar text-white bg-dark h-100 p-0">
				<Container>
					<h1>AI Dashboard</h1>
				</Container>
				<ListGroup className="list-group-flush">
					<ListGroupItem className="bg-dark border-bottom">
						<Container className="text-end">
							<Button onClick={handleAddNewBotButtonClick}>
								<FontAwesomeIcon icon={faPlus} className="me-2" />
								Add New Bot
							</Button>
						</Container>
					</ListGroupItem>
					<ListGroupItem className="bg-dark text-white border-bottom">
						<Container className="text-start">
							Bot Name
						</Container>
					</ListGroupItem>
				</ListGroup>
			</Col>
			<Col className="col-8 h-100">
				<Container>
					<CreateBotForm />
				</Container>
			</Col>
		</Row>
	);
}

function handleAddNewBotButtonClick() {
	console.log("Add new bot button clicked");
}

function CreateBotForm() {
	const [botData, setBotData] = useState({
		name: "",
		description: "",
		model: "gpt-4",
		personality: "",
		user_history: ""
	});
	
	let handleChange = (event) => {
		setBotData({
			...botData,
			[event.target.name]: event.target.value
		});
	}
	
	let handleSubmit = (event) => {
		event.preventDefault();
		
		console.log(botData);	
	}
	
	return (
		<div>
			<h1>Create New Bot</h1>
			<form onSubmit={handleSubmit}>
				<div className="mb-3">
					<label htmlFor="name" className="form-label">Bot Name*</label>
					<input onChange={handleChange} type="text" className="form-control" id="create-bot-form-name" name="name" placeholder="Enter bot name" required value={botData.name} />
				</div>
				<div className="mb-3">
					<label htmlFor="description" className="form-label">Description</label>
					<input onChange={handleChange} type="text" className="form-control" id="create-bot-form-description" name="description" placeholder="Enter bot description" value={botData.description} />
				</div>
				<div className="mb-3">
					<label htmlFor="model" className="form-label">Model*</label>
					<select onChange={handleChange} id="create-bot-form-model" name="model" className="form-control" value={botData.model}>
						<option value="gpt-4">GPT-4</option>
					</select>
				</div>
				<div className="mb-3">
					<label htmlFor="personality" className="form-label">Personality</label>
					<textarea onChange={handleChange} className="form-control" id="create-bot-form-personality" name="personality" rows="3" placeholder="Enter bot personality" value={botData.personality} ></textarea>
				</div>
				<div className="mb-3">
					<label htmlFor="user_history" className="form-label">User History</label>
					<textarea onChange={handleChange} className="form-control" id="create-bot-form-user_history" name="user_history" rows="3" placeholder="Enter user history" value={botData.user_history} ></textarea>
				</div>
				<button type="submit" className="btn btn-primary">Create Bot</button>
			</form>
		</div>
	);
}

export default App;
