import './App.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faComments, faGear, faPlus, faRobot } from '@fortawesome/free-solid-svg-icons';
import { useDispatch, useSelector } from "react-redux";
import { PAGE_STATUSES, goToBotListPage, goToChatListPage, goToChatEditPage, selectPageStatus, selectSelectedChat, goToChatPage } from './store/page';
import ChatList from './components/ChatList';
import Chat from './components/Chat';
import ChatForm from './forms/ChatForm';
import CreateBotForm from './forms/BotForm';
import { useState } from 'react';
import { debounce } from 'lodash';
import { persistChat } from './store/chats';
import BotList from './components/BotList';

function App() {
	const dispatch = useDispatch();
	const pageStatus = useSelector(selectPageStatus);
	const selectedChat = useSelector(selectSelectedChat);
	
	let title = "AI Dashboard";
	
	let content;
	switch(pageStatus) {
		case PAGE_STATUSES.CREATE_BOT:
			content = <CreateBotForm/>;
		break;
		case PAGE_STATUSES.BOT_LIST:
			content = <BotList/>;
		break;
		case PAGE_STATUSES.BOT_CHAT:
			title = <ChatTitle/>;
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
	
	let clickChatEditButton = () => {
		if (pageStatus === PAGE_STATUSES.CREATE_CHAT) {
			dispatch(goToChatPage(selectedChat));
		} else {
			dispatch(goToChatEditPage(selectedChat));
		}
	}
	
	let chatEditButton = "";
	if (selectedChat) {
		chatEditButton = (
			<span className="btn bg-white border mx-1" onClick={clickChatEditButton}>
				<FontAwesomeIcon icon={faGear}/>
			</span>
		);
	}
	
	return (
		<div className="App h-100 d-flex flex-column">
			<nav className="navbar bg-light border-bottom py-0">
				<div className="mx-2 d-flex w-100 align-items-center py-2">
					<span className="btn bg-white border" onClick={() => dispatch(goToBotListPage())}>
						<FontAwesomeIcon icon={faRobot}/>
					</span>
					<span className="btn bg-white border mx-1" onClick={() => dispatch(goToChatListPage(selectedChat))}>
						<FontAwesomeIcon icon={faComments}/>
					</span>
					
					<div className="flex-grow-1">
						{title}
					</div>
					
					<div className="d-flex">
						{chatEditButton}
						<span className="btn bg-white border" onClick={() => dispatch(goToChatEditPage())}>
							<FontAwesomeIcon icon={faPlus}/>
						</span>
					</div>
				</div>
			</nav>
			<div className="h-100 d-flex flex-column">
				{content}
			</div>
		</div>
	);
}

const debounceSave = debounce((dispatch, newChat) => {
	dispatch(persistChat(newChat));
}, 1000);

function ChatTitle() {
	const dispatch = useDispatch();
	const selectedChat = useSelector(selectSelectedChat);
	
	const [formData, setFormData] = useState({
		name: selectedChat.name,
	});
	
	const handleChange = (event) => {
		let newFormData = {
			...formData,
			[event.target.name]: event.target.value
		};
		setFormData(newFormData);
		let newChat = Object.assign({}, selectedChat, newFormData);
		debounceSave(dispatch, newChat);
	}
	
	return (
		<input 
			className="bg-light flex-grow-1 form-control border-0 text-center" 
			name="name" 
			value={formData.name || ''} 
			onChange={handleChange}
		/>
	);
}

export default App;
