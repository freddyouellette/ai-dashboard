import { useDispatch, useSelector } from "react-redux";
import { selectSelectedChat } from "../store/page";
import { useState, useEffect, useRef } from "react"; // import useRef
import { Button } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheck, faCommentDots, faCopy, faEllipsisH, faEllipsisV, faEyeSlash, faPaperPlane, faExclamationTriangle } from "@fortawesome/free-solid-svg-icons";
import { getChatMessages, selectMessages, selectWaitingForCorrectionId, selectWaitingForResponse, sendMessage, updateMessageActive, updateMessageBreakAfter } from "../store/messages";
import './chat.css'
import { getBots, selectBots } from "../store/bots";
import moment from "moment";
import popup from "../util/popup"
import { diffWords } from "diff";
import MessageText from "./MessageText";

export default function Chat() {
	const dispatch = useDispatch();
	const selectedChat = useSelector(selectSelectedChat);
	const { messages, messagesLoading, messagesError, responseFailed } = useSelector(selectMessages);
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
	
	const dispatchSendMessage = (role, getResponse) => {
		if (messageToSend === "") {
			if (getResponse && role != "BOT") {
				return popup.confirm("The message field is empty. Are you sure you would like to send only previous messages?").then(ok => {
					if (ok) dispatch(sendMessage(selectedChat.ID, selectedChat?.bot_id, messageToSend, role, getResponse));
				});
			} else {
				// do nothing
			}
		} else {
			dispatch(sendMessage(selectedChat.ID, selectedChat?.bot_id, messageToSend, role, getResponse))
			setMessageToSend("")
		}
	}
	
	const handleMessageTextKeyDown = event => {
		if ((event.ctrlKey || event.metaKey) && (event.keyCode === 13)) {
			dispatchSendMessage("USER", true)
		}
	}

	const focusTextArea = () => {
		const textArea = document.querySelector("#chat-message");
		if (textArea) {
		  textArea.focus();
		}
	  };
	
	  focusTextArea();
	
	let messagesArr = Object.values(messages);
	let lastMessage = messagesArr[messagesArr.length - 1] || null;
	let lastMessageBreakAfter = ""
	if (lastMessage) {
		if (!lastMessage.break_after) {
			lastMessageBreakAfter = <li><div className="dropdown-item cursor-pointer" onClick={() => dispatch(updateMessageBreakAfter(lastMessage, true))}>Add Break</div></li>
		} else {
			lastMessageBreakAfter = <li><div className="dropdown-item cursor-pointer" onClick={() => dispatch(updateMessageBreakAfter(lastMessage, false))}>Remove Break</div></li>
		}
	}
	return (
		<div className="d-flex flex-column flex-grow-1">
			<div className="message-list flex-grow-1 overflow-auto mb-2 border-bottom px-2" style={{height: '0px'}} ref={messageListRef}>
				<div>
					<SystemMessage title="Bot Personality:" content={chatBot.personality} />
					<SystemMessage title="User History:" content={chatBot.user_history} />
					<SystemMessage title="Chat Instructions:" content={selectedChat.instructions} />
					{messagesArr.map((message, i) => {
						let barrier = "";
						const memory = moment().subtract(selectedChat.memory_duration, 'seconds');
						if (message.break_after) {
							barrier = <div className="hr-with-text my-3">
								<em className="text-muted" style={{fontSize: "0.7em"}}>Break</em>
							</div>
						} else if (
							moment(message.CreatedAt).isBefore(memory) 
							&& (
								!messagesArr[i+1]
								|| moment(messagesArr[i+1].CreatedAt).isAfter(memory)
							)
						) {
							barrier = <div className="hr-with-text my-3">
								<em className="text-muted" style={{fontSize: "0.7em"}}>Memory Limit</em>
							</div>
						}
						switch (message.role) {
							case "USER":
								let messageDiv = <MessageText>{message.text}</MessageText>;
								if (message.correction) {
									let messageParts = [];
									diffWords(message.text, message.correction).forEach((part, i) => {
										if (part.added) {
											messageParts.push(<span key={i} className="text-success fw-bold">{part.value}</span>);
										} else if (part.removed) {
											messageParts.push(<span key={i} className="text-danger text-decoration-line-through fst-italic me-1">{part.value}</span>);
										} else {
											messageParts.push(<span key={i}>{part.value}</span>);
										}
									});
									messageDiv = <div className="text-break">{messageParts}</div>;
								}
								return (
									<div key={message.ID}>
										<div className={"text-start user-message p-2 m-2 ms-5 rounded " + (message.active ? "active" : "inactive")}>
											<div className="d-flex justify-content-between align-items-center">
												<div><b className="">👤 You:</b></div>
												<div className="d-flex gap-1 align-items-center">
													{message.ID === waitingForCorrectionId ? <div className="help-text text-danger">Correcting...</div> : ""}
													{message.correction ? <FontAwesomeIcon className="text-success" icon={faCheck} /> : ""}
													<div className="help-text">{moment(message.CreatedAt).fromNow()}</div>
													<MessageContextMenu message={message} />
												</div>
											</div>
											{messageDiv}
										</div>
										{barrier}
									</div>
								)
							case "BOT":
								return (
									<div key={message.ID}>
										<div className={"text-start bot-message p-2 m-2 me-5 rounded " + (message.active ? "active" : "inactive")}>
											<div className="d-flex justify-content-between align-items-center">
												<div><b>🤖 {chatBot.name}:</b></div>
												<div className="d-flex align-items-start">
													<div className="help-text mx-1">{moment(message.CreatedAt).fromNow()}</div>
													<MessageContextMenu message={message} />
												</div>
											</div>
											<MessageText>{message.text}</MessageText>
										</div>
										{barrier}
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
					{responseFailed ? (
						<div className={"text-start bot-message p-2 m-2 me-5 rounded active border border-danger"}>
							<div className="d-flex justify-content-between align-items-center">
								<div><b>🤖 {chatBot.name}:</b></div>
							</div>
							<div className="text-danger"><FontAwesomeIcon icon={faExclamationTriangle} /> Failed to get response. Please try again</div>
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
					<div className="btn-group-vertical ms-2">
						<Button className="py-3 border-bottom" onClick={() => dispatchSendMessage("USER", true)}><FontAwesomeIcon icon={faPaperPlane}/></Button>
						<Button className="py-0 border-top rounded-bottom" data-bs-toggle="dropdown"><FontAwesomeIcon icon={faEllipsisH}/></Button>
						<ul className="dropdown-menu">
							<li><div className="dropdown-item cursor-pointer" onClick={() => dispatchSendMessage("USER", false)}>Send Without Response</div></li>
							<li><div className="dropdown-item cursor-pointer" onClick={() => dispatchSendMessage("BOT", false)}>Create Bot Message</div></li>
							{lastMessageBreakAfter}
						</ul>
					</div>
				</div>
				<div className="text-start ms-2 d-none d-md-flex">
					<small><em>Ctrl+Enter / ⌘+Enter to send message</em></small>
				</div>
			</div>
		</div>
	);
}

function SystemMessage({ title, content }) {
	if (!content) return "";
	return (
		<div className="text-start system-message p-2 m-2 rounded border text-muted">
			<div className="d-flex justify-content-between">
				<b className="">{title}</b>
				<CopyButton text={content} />
			</div>
			<MessageText>{content}</MessageText>
		</div>
	);
}

const copyToClipboard = (setCopied, text) => {
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

function MessageContextMenu({ message }) {
	const dispatch = useDispatch();
	
	return (
		<>
			{message.active ? "" : <FontAwesomeIcon className="text-muted" icon={faEyeSlash} size="sm"/>}
			<div className="btn-group dropstart cursor-pointer d-flex align-items-center">
				<FontAwesomeIcon className="dropdown-toggle text-muted ps-2 pe-1" data-bs-toggle="dropdown" icon={faEllipsisV} size="lg"/>
				<ul className="dropdown-menu">
					<div className="dropdown-item" onClick={() => copyToClipboard(() => {}, message.text)}>Copy Text</div>
					{
						message.active 
						? <div className="dropdown-item" onClick={() => dispatch(updateMessageActive(message, false))}>Hide from Bot</div>
						: <div className="dropdown-item" onClick={() => dispatch(updateMessageActive(message, true))}>Show to Bot</div>
					}
					{
						message.break_after 
						? <div className="dropdown-item" onClick={() => dispatch(updateMessageBreakAfter(message, false))}>Remove Break After</div>
						: <div className="dropdown-item" onClick={() => dispatch(updateMessageBreakAfter(message, true))}>Add Break After</div>
					}
				</ul>
			</div>
		</>
	);
}

function CopyButton({ text }) {
	const [copied, setCopied] = useState(false);
	
	return (
		<FontAwesomeIcon className={copied ? "copy-button text-success" : "copy-button"} icon={faCopy} onClick={() => copyToClipboard(setCopied, text)} />
	)
}