import { Col, Container, Row } from 'react-bootstrap';
import './App.css';
import BotForm from './forms/BotForm';
import Chat from './components/Chat';
import { useSelector } from 'react-redux'
import { selectPageStatus, PAGE_STATUSES } from './store/page'
import SidebarMenu from './components/SidebarMenu';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faBars, faGear, faGripLines } from '@fortawesome/free-solid-svg-icons';

function App() {
	// const chatSelected = useSelector(selectChatSelected);
	const chatSelected = {};
	// const botSelected = useSelector(selectBotSelected);
	const botSelected = {
		name: "Uncle Bob",
		description: "",
		model: "",
		personality: "",
		user_history: "",
		randomness: "",
	};
	
	
	return (
		<div className="App h-100">
			<nav className="navbar px-2 bg-light border-bottom">
				<div data-bs-toggle="collapse" data-bs-target="#navbar-menu">
					<span className="btn border"><FontAwesomeIcon icon={faBars}/></span>
				</div>
				
				<div className="flex-grow-1">
					{botSelected.name}
				</div>
				
				<div>
					<span className="btn border"><FontAwesomeIcon icon={faGear}/></span>
				</div>
				
				<div className="collapse navbar-collapse" id="navbar-menu">
					<ul className="navbar-nav me-auto">
						<li className="nav-item active">
							<a className="nav-link" href="#">Home</a>
						</li>
					</ul>
				</div>
			</nav>
		</div>
	);
}

export default App;
