import { useDispatch, useSelector } from "react-redux";
import { selectSelectedChat } from "../store/page";
import { useState, useEffect, useRef } from "react"; // import useRef
import { Button } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCommentDots, faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { getChatMessages, selectMessages, selectWaitingForResponse, sendMessage } from "../store/messages";
import Markdown from 'react-markdown'
import './chat.css'
import { getBots, selectBots } from "../store/bots";

export default function Chat() {
	const dispatch = useDispatch();
	const selectedChat = useSelector(selectSelectedChat);
	const { messages, messagesLoading, messagesError } = useSelector(selectMessages);
	const waitingForResponse = useSelector(selectWaitingForResponse);
	const [messageToSend, setMessageToSend] = useState("");
	const { bots, botsLoading, botsError } = useSelector(selectBots);
	const chatBot = bots[selectedChat?.bot_id];
	
	useEffect(() => {
		dispatch(getBots());
	}, [dispatch]);
	
	useEffect(() => {
		dispatch(getChatMessages(selectedChat))
	}, [selectedChat, dispatch]);
	
	// scroll to bottom of list when a new message appears
	const messageListRef = useRef(null);
	useEffect(() => {
		if (messageListRef.current !== null) {
			messageListRef.current.scrollTop = messageListRef.current.scrollHeight;
		}
	}, [messages]);
	
	if (messagesLoading || botsLoading) return <div>Loading...</div>;
	if (messagesError) return <div>Error loading messages...</div>;
	if (botsError) return <div>Error loading bots...</div>;
	
	let handleMessageToSendChange = event => {
		setMessageToSend(event.target.value)
	}
	
	let dispatchSendMessage = () => {
		dispatch(sendMessage(selectedChat.ID, messageToSend))
		setMessageToSend("")
	}
	
	let handleMessageTextKeyDown = event => {
		if ((event.ctrlKey || event.metaKey) && (event.keyCode === 13)) {
			dispatchSendMessage()
		}
	}

	return (
		<div className="d-flex flex-column flex-grow-1">
			<div className="message-list flex-grow-1 overflow-auto mb-2 border-bottom px-2" style={{height: '0px'}} ref={messageListRef}>
				<div>
					<div className="text-start system-message p-2 m-2 rounded border text-muted">
						<b className="">Bot Personality:</b>
						<Markdown>
							{chatBot.personality}
						</Markdown>
					</div>
					{Object.values(messages).map(message => {
						switch (message.role) {
							case "USER":
								return (
									<div key={message.ID} className="text-start user-message p-2 m-2 ms-5 rounded border">
										<b className="">👤 You:</b>
										<Markdown>
											{message.text}
										</Markdown>
									</div>
								)
							case "BOT":
								return (
									<div key={message.ID} className="text-start bg-light p-2 m-2 me-5 rounded border">
										<b>🤖 {chatBot.name}:</b>
										<Markdown>
											{message.text}
										</Markdown>
									</div>
								)
							default:
								return ""
						}
					})}
					{waitingForResponse ? (
						<div className="d-flex justify-content-start">
							<div className="m-2 p-3 border rounded bg-light">
								🤖 <FontAwesomeIcon className="ms-1 mb-2" icon={faCommentDots} bounce />
							</div>
						</div>
					) : ""}
				</div>
			</div>
			<div className="mb-2 px-2">
				<div className="d-flex">
					<textarea 
						onChange={handleMessageToSendChange} 
						onKeyDown={handleMessageTextKeyDown}
						type="text" 
						className="form-control" 
						id="chat-message" 
						name="message" 
						placeholder="Enter message" 
						required 
						value={messageToSend}></textarea>
					<Button className="ms-3" onClick={dispatchSendMessage}><FontAwesomeIcon icon={faPaperPlane} /></Button>
				</div>
				<div className="text-start ms-2 d-none d-md-flex">
					<small><em>Ctrl+Enter / ⌘+Enter to send message</em></small>
				</div>
			</div>
		</div>
	);
}