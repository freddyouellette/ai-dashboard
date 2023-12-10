import './App.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faBars, faGear } from '@fortawesome/free-solid-svg-icons';
import ChatList from './components/ChatList';
import { TODO } from './util/todo';

function App() {
	// const chatSelected = useSelector(selectChatSelected);
	// const chatSelected = {};
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
					<span className="btn border" onClick={TODO}><FontAwesomeIcon icon={faGear}/></span>
				</div>
				
				<div className="collapse navbar-collapse" id="navbar-menu">
					<ChatList/>
				</div>
			</nav>
		</div>
	);
}

export default App;
