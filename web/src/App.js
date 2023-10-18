import { Button, Col, Container, ListGroup, ListGroupItem, Row } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faEdit, faPlus } from '@fortawesome/free-solid-svg-icons'
import './App.css';
import { useState, useEffect } from 'react';
import axios from 'axios';
import BotForm from './forms/BotForm';

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
			content = <BotForm upsertBotFunc={upsertBot} botToUpdate={botToUpdate} />;
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

export default App;
