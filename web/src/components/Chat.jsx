import { useDispatch, useSelector } from "react-redux";
import { selectSelectedChat } from "../store/page";
import { useState, useEffect, useRef } from "react"; // import useRef
import { Button } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck, faCommentDots, faCopy, faPaperPlane} from "@fortawesome/free-solid-svg-icons";
import { getChatMessages, selectMessages, selectWaitingForCorrectionId, selectWaitingForResponse, sendMessage } from "../store/messages";
import './chat.css'
import { getBots, selectBots } from "../store/bots";
import moment from "moment";
import popup from "../util/popup"
import { diffWords } from "diff";
import MessageText from "./MessageText";

export default function Chat() {
	const dispatch = useDispatch();
	const selectedChat = useSelector(selectSelectedChat);
	const { messages, messagesLoading, messagesError } = useSelector(selectMessages);
	const waitingForResponse = useSelector(selectWaitingForResponse);
	const [messageToSend, setMessageToSend] = useState("");
	const { bots, botsLoading, botsError } = useSelector(selectBots);
	const chatBot = bots[selectedChat?.bot_id];
	
	const waitingForCorrectionId = useSelector(selectWaitingForCorrectionId);
	
	useEffect(() => {
		dispatch(getBots());
	}, [dispatch]);
	
	useEffect(() => {
		dispatch(getChatMessages(selectedChat));
	}, [selectedChat, dispatch]);
	
	// scroll to bottom of list when a new message appears
	const messageListRef = useRef(null);
	useEffect(() => {
		if (messageListRef.current !== null) {
			messageListRef.current.scrollTop = messageListRef.current.scrollHeight;
		}
	}, [messages, bots]);
	
	if (messagesLoading || botsLoading) return <div>Loading...</div>;
	if (messagesError) return <div>Error loading messages...</div>;
	if (botsError) return <div>Error loading bots...</div>;
	if (!chatBot) return <div className="pt-3 text-danger">Unknown Bot...</div>;
	
	const handleMessageToSendChange = event => {
		setMessageToSend(event.target.value)
	}
	
	const dispatchSendMessage = () => {
		if (messageToSend === "") {
			return popup.confirm("The message field is empty. Are you sure you would like to send only previous messages?").then(ok => {
				if (ok) dispatch(sendMessage(selectedChat.ID, selectedChat?.bot_id, messageToSend));
			});
		} else {
			dispatch(sendMessage(selectedChat.ID, selectedChat?.bot_id, messageToSend))
			setMessageToSend("")
		}
	}
	
	const handleMessageTextKeyDown = event => {
		if ((event.ctrlKey || event.metaKey) && (event.keyCode === 13)) {
			dispatchSendMessage()
		}
	}

	const focusTextArea = () => {
		const textArea = document.querySelector("#chat-message");
		if (textArea) {
		  textArea.focus();
		}
	  };
	
	  focusTextArea();
	
	let personalityMessage = "";
	if (chatBot.personality) {
		personalityMessage = (
			<div className="text-start system-message p-2 m-2 rounded border text-muted">
				<div className="d-flex justify-content-between">
					<b className="">Bot Personality:</b>
					<CopyButton text={chatBot.personality} />
				</div>
				<MessageText>{chatBot.personality}</MessageText>
			</div>
		)
	}
	
	let userHistoryMessage = "";
	if (chatBot.user_history) {
		userHistoryMessage = (
			<div className="text-start system-message p-2 m-2 rounded border text-muted">
				<div className="d-flex justify-content-between">
					<b className="">User History:</b>
					<CopyButton text={chatBot.user_history} />
				</div>
				<MessageText>{chatBot.user_history}</MessageText>
			</div>
		)
	}
	
	return (
		<div className="d-flex flex-column flex-grow-1">
			<div className="message-list flex-grow-1 overflow-auto mb-2 border-bottom px-2" style={{height: '0px'}} ref={messageListRef}>
				<div>
					{personalityMessage}
					{userHistoryMessage}
					{Object.values(messages).map(message => {
						switch (message.role) {
							case "USER":
								let messageDiv = <MessageText>{message.text}</MessageText>;
								if (message.correction) {
									let messageParts = [];
									diffWords(message.text, message.correction).forEach(part => {
										if (part.added) {
											messageParts.push(<span class="text-success fw-bold">{part.value}</span>);
										} else if (part.removed) {
											messageParts.push(<span class="text-danger text-decoration-line-through fst-italic me-1">{part.value}</span>);
										} else {
											messageParts.push(<span>{part.value}</span>);
										}
									});
									messageDiv = <div className="text-break">{messageParts}</div>;
								}
								return (
									<div key={message.ID} className="text-start user-message p-2 m-2 ms-5 rounded border">
										<div className="d-flex justify-content-between align-items-center">
											<div><b className="">👤 You:</b></div>
											<div className="d-flex gap-1 align-items-center">
												{message.ID === waitingForCorrectionId ? <div className="help-text text-danger">Correcting...</div> : ""}
												{message.correction ? <FontAwesomeIcon className="text-success" icon={faCheck} /> : ""}
												<div className="help-text">{moment(message.CreatedAt).fromNow()}</div>
												<CopyButton text={message.text} />
											</div>
										</div>
										{messageDiv}
									</div>
								)
							case "BOT":
								return (
									<div key={message.ID} className="text-start bg-light p-2 m-2 me-5 rounded border">
										<div className="d-flex justify-content-between align-items-center">
											<div><b>🤖 {chatBot.name}:</b></div>
											<div className="d-flex">
												<div className="help-text mx-1">{moment(message.CreatedAt).fromNow()}</div>
												<CopyButton text={message.text} />
											</div>
										</div>
										<MessageText>{message.text}</MessageText>
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

function CopyButton({ text }) {
	const [copied, setCopied] = useState(false);
	
	const copyToClipboard = (text) => {
		if (!navigator.clipboard) {
			console.error('Clipboard API not available');
			return;
		}
		navigator.clipboard.writeText(text).then(() => {
			setCopied(true);
			setTimeout(() => setCopied(false), 2000);
		}).catch(err => {
			console.error('Failed to copy text: ', err);
		});
	}
	
	// if (window.isSecureContext === false) return "";
	
	return (
		<FontAwesomeIcon className={copied ? "copy-button text-success" : "copy-button"} icon={faCopy} onClick={() => copyToClipboard(text)} />
	)
}