import './App.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faBars, faGear, faPlus } from '@fortawesome/free-solid-svg-icons';
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
			<nav className="navbar bg-light border-bottom py-0">
				<div class="mx-2 d-flex w-100 align-items-center py-2">
					<div data-bs-toggle="collapse" data-bs-target="#navbar-menu">
						<span className="btn border"><FontAwesomeIcon icon={faBars}/></span>
					</div>
					
					<div className="flex-grow-1">
						AI Dashboard
					</div>
					
					<div data-bs-toggle="collapse" data-bs-target="#navbar-menu">
						<span className="btn border"><FontAwesomeIcon icon={faPlus}/></span>
					</div>
				</div>
				
				<div className="collapse navbar-collapse border-top" id="navbar-menu">
					<div className="border-bottom py-3">Bots</div>
					<div className="py-3">Chats</div>
				</div>
			</nav>
		</div>
	);
}

export default App;
