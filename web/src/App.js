import { Col, Container, Row } from 'react-bootstrap';
import './App.css';
import BotForm from './forms/BotForm';
import Chat from './components/Chat';
import { useSelector } from 'react-redux'
import { selectPageStatus, PAGE_STATUSES } from './store/page'
import SidebarMenu from './components/SidebarMenu';

function App() {
	const pageStatus = useSelector(selectPageStatus)
	
	let content;
	switch(pageStatus) {
		case PAGE_STATUSES.CREATE_BOT:
			content = <BotForm />;
		break;
		case PAGE_STATUSES.BOT_CHAT:
			content = <Chat />
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
				<SidebarMenu />
			</Col>
			<Col className="col-8 h-100 d-flex flex-column">
				<Container className="flex-grow-1 d-flex flex-column">
					{content}
				</Container>
			</Col>
		</Row>
	);
}

export default App;
