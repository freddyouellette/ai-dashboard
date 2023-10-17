import { Button, Col, Container, ListGroup, ListGroupItem, Row } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faEdit, faPlus } from '@fortawesome/free-solid-svg-icons'
import './App.css';
import { useState, useEffect } from 'react';
import axios from 'axios';

const PAGE_CONTENT_CREATE_BOT = "create-bot";

function App() {
	const [bots, setBots] = useState([])
	const [pageContent, setPageContent] = useState(null)
	let [botToUpdate, setBotToUpdate] = useState(null)
	let setBotToUpdateAndRender = (bot) => {
		setBotToUpdate(bot);
		setPageContent(PAGE_CONTENT_CREATE_BOT);
	}
	
	useEffect(() => {
		axios.get("http://localhost:8080/bots")
		.then(response => {
			console.log("Bots", response.data);
			setBots(response.data);
		}).catch(error => {
			console.error(error);
		})
	}, [])
	
	let upsertBot = (bot) => {
		if (bot.ID) {
			// Update bot
			let botIndex = bots.findIndex(b => b.ID === bot.ID);
			if (botIndex !== -1) {
				setBots([
					...bots.slice(0, botIndex),
					bot,
					...bots.slice(botIndex + 1)
				]);
			}
		} else {
			// Create bot
			setBots([
				...bots,
				bot
			]);
		}
	}
	
	let content;
	switch(pageContent) {
		case PAGE_CONTENT_CREATE_BOT:
			content = <CreateBotForm upsertBotFunc={upsertBot} botToUpdate={botToUpdate} />;
		break;
		default:
			content = <div className="text-center text-italics">Select a bot</div>;
	}
	
	return (
		<Row className="App h-100">
			<Col className="sidebar text-white bg-dark h-100 p-0">
				<Container>
					<h1>AI Dashboard</h1>
				</Container>
				<ListGroup className="list-group-flush">
					<ListGroupItem className="bg-dark border-bottom">
						<Container className="text-center">
							<Button onClick={() => {setBotToUpdateAndRender(null)}}>
								<FontAwesomeIcon icon={faPlus} className="me-2" />
								Add New Bot
							</Button>
						</Container>
					</ListGroupItem>
					{bots.map(bot => {
						return (
							<ListGroupItem key={bot.ID} className="bg-dark text-white border-bottom">
								<Container className="text-start">
									<div className="d-flex justify-content-between align-items-center">
										<strong>{bot.name}</strong>
										<FontAwesomeIcon icon={faEdit} className="ms-2 cursor-pointer" style={{"cursor": "pointer"}} onClick={() => {setBotToUpdateAndRender(bot)}} />
									</div>
									<div>{bot.description}</div>
								</Container>
							</ListGroupItem>
						);
					})}
				</ListGroup>
			</Col>
			<Col className="col-8 h-100">
				<Container>
					{content}
				</Container>
			</Col>
		</Row>
	);
}

function CreateBotForm({upsertBotFunc, botToUpdate}) {
	const [botFormData, setBotFormData] = useState(null);
	
	// update the state if the botToUpdate prop changes
	useEffect(() => {
		setBotFormData(botToUpdate);
	}, [botToUpdate])
	
	let handleChange = (event) => {
		setBotFormData({
			...botFormData,
			[event.target.name]: event.target.value
		});
	}
	
	let handleSubmit = async (event) => {
		event.preventDefault();
		
		axios[botFormData.ID ? "put" : "post"]("http://localhost:8080/bots", botFormData)
		.then(response => {
			console.log(response);
			upsertBotFunc(response.data);
		}).catch(error => {
			console.error(error);
		})
	}
	
	return (
		<div>
			<h1>{botFormData?.ID ? 'Update Bot' : 'Create New Bot'}</h1>
			<form onSubmit={handleSubmit}>
				<div className="mb-3">
					<label htmlFor="name" className="form-label">Bot Name*</label>
					<input onChange={handleChange} type="text" className="form-control" id="create-bot-form-name" name="name" placeholder="Enter bot name" required value={botFormData?.name || ''} />
				</div>
				<div className="mb-3">
					<label htmlFor="description" className="form-label">Description</label>
					<input onChange={handleChange} type="text" className="form-control" id="create-bot-form-description" name="description" placeholder="Enter bot description" value={botFormData?.description || ''} />
				</div>
				<div className="mb-3">
					<label htmlFor="model" className="form-label">Model*</label>
					<select onChange={handleChange} id="create-bot-form-model" name="model" className="form-control" value={botFormData?.model ||  ''}>
						<option value="gpt-4">GPT-4</option>
					</select>
				</div>
				<div className="mb-3">
					<label htmlFor="personality" className="form-label">Personality</label>
					<textarea onChange={handleChange} className="form-control" id="create-bot-form-personality" name="personality" rows="3" placeholder="Enter bot personality" value={botFormData?.personality ||  ''} ></textarea>
				</div>
				<div className="mb-3">
					<label htmlFor="user_history" className="form-label">User History</label>
					<textarea onChange={handleChange} className="form-control" id="create-bot-form-user_history" name="user_history" rows="3" placeholder="Enter user history" value={botFormData?.user_history ||  ''} ></textarea>
				</div>
				<Button className="btn-secondary me-3">Cancel</Button>
				<button type="submit" className="btn btn-primary">{botFormData?.ID ? 'Update' : 'Create'} Bot</button>
			</form>
		</div>
	);
}

export default App;
