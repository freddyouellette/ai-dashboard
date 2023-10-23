import { useDispatch, useSelector } from "react-redux";
import { selectChatBot, selectSelectedChat } from "../store/page";
import { useState } from "react";
import { Button } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { selectMessages, sendMessage } from "../store/messages";
import Markdown from 'react-markdown'

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
	
	return (
		<div className="bg-red flex-grow-1 d-flex flex-column">
			<div className="flex-grow-1">
				<div className="overflow-auto">
					{messages.map(message => {
						switch (message.role) {
							case "USER":
								return (
									<div key={message.ID} className="text-start bg-light p-2 m-2 ms-5 rounded border">
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
			<div className="mb-3">
				<div className="d-flex">
					<textarea onChange={handleMessageToSendChange}  type="text" className="form-control" id="chat-message" name="message" placeholder="Enter message" required value={messageToSend}></textarea>
					<Button className="ms-3" onClick={() => dispatch(sendMessage(selectedChat.ID, messageToSend))}><FontAwesomeIcon icon={faPaperPlane} /></Button>
				</div>
				<div className="text-start ms-2">
					<small><em>Ctrl+Enter / ⌘+Enter to send message</em></small>
				</div>
			</div>
		</div>
	);
}