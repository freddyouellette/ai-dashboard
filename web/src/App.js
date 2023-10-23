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
		<div className="App h-100 container-fluid text-start">
			<Row className="h-100">
				<Col className="col-4 sidebar text-white bg-dark h-100 p-0">
					<Container>
						<h1>AI Dashboard</h1>
					</Container>
					<SidebarMenu />
				</Col>
				<Col className="col-8 h-100 d-flex flex-column p-0">
					<div className="flex-grow-1 d-flex flex-column">
						{content}
					</div>
				</Col>
			</Row>
		</div>
	);
}

export default App;
