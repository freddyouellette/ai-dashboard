import './App.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faBars, faPlus } from '@fortawesome/free-solid-svg-icons';
import { useDispatch, useSelector } from "react-redux";
import { PAGE_STATUSES, goToBotListPage, goToChatListPage, goToCreateChatPage, selectPageStatus } from './store/page';
import ChatList from './components/ChatList';
import Chat from './components/Chat';
import ChatForm from './forms/ChatForm';

function App() {
	const dispatch = useDispatch();
	const pageStatus = useSelector(selectPageStatus);
	
	let content;
	switch(pageStatus) {
		case PAGE_STATUSES.CREATE_BOT:
			content = "Create Bot";
		break;
		case PAGE_STATUSES.BOT_LIST:
			content = "Bot List";
		break;
		case PAGE_STATUSES.BOT_CHAT:
			content = <Chat/>;
		break;
		case PAGE_STATUSES.CREATE_CHAT:
			content = <ChatForm/>;
		break;
		case PAGE_STATUSES.CHAT_LIST:
			content = <ChatList/>;
		break;
		default:
			content = "";
		break;
	}
	
	return (
		<div className="App h-100 d-flex flex-column">
			<nav className="navbar bg-light border-bottom py-0">
				<div className="mx-2 d-flex w-100 align-items-center py-2">
					<div data-bs-toggle="collapse" data-bs-target="#navbar-menu">
						<span className="btn border">
							<FontAwesomeIcon icon={faBars}/>
						</span>
					</div>
					
					<div className="flex-grow-1">
						AI Dashboard
					</div>
					
					<div>
						<span className="btn border" onClick={() => dispatch(goToCreateChatPage())}>
							<FontAwesomeIcon icon={faPlus}/>
						</span>
					</div>
				</div>
				
				<div className="collapse navbar-collapse" id="navbar-menu">
					<NavItem name="Bots" onClick={() => dispatch(goToBotListPage())}/>
					<NavItem name="Chats" onClick={() => dispatch(goToChatListPage())}/>
				</div>
			</nav>
			{content}
		</div>
	);
}

function NavItem({ name, onClick }) {
	return (
		<div 
			className="border-top py-3 cursor-pointer" 
			data-bs-toggle="collapse" 
			data-bs-target="#navbar-menu" 
			onClick={onClick}>
			{name}
		</div>
	);
}

export default App;
