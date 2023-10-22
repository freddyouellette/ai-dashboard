import { useDispatch, useSelector } from "react-redux";
import { selectSelectedChat } from "../store/page";
import { useState } from "react";
import { Button } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";
import { selectMessages, sendMessage } from "../store/messages";

export default function Chat() {
	let dispatch = useDispatch()
	let selectedChat = useSelector(selectSelectedChat)
	let messages = useSelector(selectMessages)
	messages = messages.filter(message => message.chat_id === selectedChat.ID)
	
	let [messageToSend, setMessageToSend] = useState("")
	
	let handleMessageToSendChange = event => {
		setMessageToSend(event.target.value)
	}
	
	return (
		<div className="bg-red flex-grow-1 d-flex flex-column">
			<div className="flex-grow-1">
				<div className="overflow-auto">
					<pre className="text-start">selectedChat: {JSON.stringify(selectedChat, null, 4)}</pre>
					<pre className="text-start">messages: {JSON.stringify(messages, null, 4)}</pre>
				</div>
			</div>
			<div className="mb-3 d-flex">
				<textarea onChange={handleMessageToSendChange} type="text" className="form-control" id="chat-message" name="message" placeholder="Enter message" required value={messageToSend}></textarea>
				<Button className="ms-3" onClick={() => dispatch(sendMessage(selectedChat.ID, messageToSend))}><FontAwesomeIcon icon={faPaperPlane} /></Button>
			</div>
		</div>
	);
}