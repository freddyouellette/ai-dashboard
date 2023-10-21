import { useSelector } from "react-redux";
import { selectSelectedChat } from "../store/chats";
import { useState } from "react";
import { Button } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPaperPlane } from "@fortawesome/free-solid-svg-icons";

export default function Chat() {
	let selectedChat = useSelector(selectSelectedChat)
	
	let [messageToSend, setMessageToSend] = useState("")
	
	let handleMessageToSendChange = event => {
		setMessageToSend(event.target.value)
	}
	
	return (
		<div className="bg-red flex-grow-1 d-flex flex-column">
			<div className="flex-grow-1">
				{JSON.stringify(selectedChat)}
			</div>
			<div className="mb-3 d-flex">
				<input onChange={handleMessageToSendChange} type="text" className="form-control" id="chat-message" name="message" placeholder="Enter message" required value={messageToSend} />
				<Button className="ms-3"><FontAwesomeIcon icon={faPaperPlane} /></Button>
			</div>
		</div>
	);
}