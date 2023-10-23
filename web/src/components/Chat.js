import { useDispatch, useSelector } from "react-redux";
import { selectChatBot, selectSelectedChat } from "../store/page";
import { useState, useEffect, useRef } from "react"; // import useRef
import { Button } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { selectMessages, sendMessage } from "../store/messages";
import Markdown from 'react-markdown'
import './chat.css'

export default function Chat() {
	let dispatch = useDispatch()
	let selectedChat = useSelector(selectSelectedChat)
	let messages = useSelector(selectMessages)
	let chatBot = useSelector(selectChatBot)
	messages = messages.filter(message => message.chat_id === selectedChat.ID)
	
	let [messageToSend, setMessageToSend] = useState("")
	
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

	// scroll to bottom of list when a new message appears
	const messageListRef = useRef(null);
	useEffect(() => {
		messageListRef.current.scrollTop = messageListRef.current.scrollHeight;
	}, [messages]);

	return (
		<div className="flex-grow-1 d-flex flex-column">
			<div className="message-list flex-grow-1 overflow-auto mb-2 border-bottom px-2" style={{height: '0px'}} ref={messageListRef}>
				<div>
					{messages.map(message => {
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
				<div className="text-start ms-2">
					<small><em>Ctrl+Enter / ⌘+Enter to send message</em></small>
				</div>
			</div>
		</div>
	);
}