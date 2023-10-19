import { Button, Col, Container, ListGroup, ListGroupItem, Row } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faEdit, faPlus } from '@fortawesome/free-solid-svg-icons'
import './App.css';
import BotForm from './forms/BotForm';
import { useSelector, useDispatch } from 'react-redux'
import { selectPageStatus, setBotToUpdate, PAGE_STATUSES } from './store/page'
import { selectBots } from './store/bots';

function App() {
	const pageStatus = useSelector(selectPageStatus)
	const dispatch = useDispatch()
	const bots = useSelector(selectBots)
	
	let content;
	switch(pageStatus) {
		case PAGE_STATUSES.CREATE_BOT:
			content = <BotForm />;
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
							<Button onClick={() => {dispatch(setBotToUpdate(null))}}>
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
										<FontAwesomeIcon icon={faEdit} className="ms-2 cursor-pointer" style={{"cursor": "pointer"}} onClick={() => {dispatch(setBotToUpdate(bot))}} />
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
