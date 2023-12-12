import './App.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faBars, faPlus } from '@fortawesome/free-solid-svg-icons';
import { useDispatch, useSelector } from "react-redux";
import { PAGE_STATUSES, goToBotEditPage, goToBotListPage, goToChatListPage, goToCreateChatPage, selectPageStatus } from './store/page';
import ChatList from './components/ChatList';
import Chat from './components/Chat';
import ChatForm from './forms/ChatForm';
import CreateBotForm from './forms/BotForm';

function App() {
	const dispatch = useDispatch();
	const pageStatus = useSelector(selectPageStatus);
	
	let content;
	switch(pageStatus) {
		case PAGE_STATUSES.CREATE_BOT:
			content = <CreateBotForm/>;
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
						<span className="btn bg-white border">
							<FontAwesomeIcon icon={faBars}/>
						</span>
					</div>
					
					<div className="flex-grow-1">
						AI Dashboard
					</div>
					
					<div>
						<span className="btn bg-white border" onClick={() => dispatch(goToCreateChatPage())}>
							<FontAwesomeIcon icon={faPlus}/>
						</span>
					</div>
				</div>
				
				<div className="collapse navbar-collapse" id="navbar-menu">
					<div 
						className="border-top p-3 cursor-pointer d-flex justify-content-between align-items-center"
						data-bs-toggle="collapse" 
						data-bs-target="#navbar-menu" 
						onClick={() => dispatch(goToBotListPage())}>
						Bots
						<div className="btn bg-white border" onClick={e => { e.stopPropagation(); dispatch(goToBotEditPage()) }}>
							<FontAwesomeIcon icon={faPlus}/>
						</div>
					</div>
					<div 
						className="border-top p-3 cursor-pointer d-flex justify-content-between align-items-center"
						data-bs-toggle="collapse" 
						data-bs-target="#navbar-menu" 
						onClick={() => dispatch(goToChatListPage())}>
						Chats
						<div className="btn bg-white border" onClick={e => { e.stopPropagation(); dispatch(goToCreateChatPage()) }}>
							<FontAwesomeIcon icon={faPlus}/>
						</div>
					</div>
				</div>
			</nav>
			<div className="h-100 d-flex flex-column">
				{content}
			</div>
		</div>
	);
}

export default App;
