import { useDispatch, useSelector } from "react-redux";
import { getBots, selectBots } from "../store/bots";
import RequiredStar from './RequiredStar';
import { persistChat } from "../store/chats";
import { useEffect, useState } from "react";
import { goToChatPage, selectSelectedChat } from "../store/page";

export default function ChatForm() {
	const dispatch = useDispatch();
	const { bots, botsLoading, botsError } = useSelector(selectBots);
	const selectedChat = useSelector(selectSelectedChat);
	const [formData, setFormData] = useState(selectedChat ?? {
		bot_id: null,
		name: 'New Chat',
		instructions: '',
		memory_duration: 60 * 60,
	});
	useEffect(() => {
		dispatch(getBots());
	}, [dispatch]);
	
	if (botsLoading) return <div>Loading...</div>
	if (botsError) return <div>Error loading bots...</div>
	
	const memory_durations = {
		"No memory": 0,
		"1 minute": 60,
		"5 minutes": 60 * 5,
		"15 minutes": 60 * 15,
		"30 minutes": 60 * 30,
		"1 hour": 60 * 60,
		"3 hours": 60 * 60 * 3,
		"12 hours": 60 * 60 * 12,
		"1 day": 60 * 60 * 24,
		"2 days": 60 * 60 * 24 * 2,
		"3 days": 60 * 60 * 24 * 3,
		"4 days": 60 * 60 * 24 * 4,
		"5 days": 60 * 60 * 24 * 5,
		"1 week": 60 * 60 * 24 * 7,
		"1 month": 60 * 60 * 24 * 28,
		"1 year": 60 * 60 * 24 * 365,
	}
	
	const handleSubmit = async (event) => {
		event.preventDefault();
		
		let createChatData = Object.assign({}, formData);
		createChatData.bot_id = parseInt(createChatData.bot_id);
		let newChat = await dispatch(persistChat(createChatData));
		dispatch(goToChatPage(newChat));
	}
	
	const handleChange = (event) => {
		setFormData({
			...formData,
			[event.target.name]: event.target.value
		});
	}
	
	return (
		<div className="container">
			<h1 className="text-center">{selectedChat ? "Edit" : "Create"} Chat</h1>
			<form onSubmit={handleSubmit}>
				<div className="mb-3">
					<label htmlFor="name" className="form-label">Name of Chat</label>
					<input onChange={handleChange} type="text" name="name" className="form-control" id="create-chat-form-name" value={formData?.name ?? ''}/>
				</div>
				<div className="mb-3">
					<label htmlFor="instructions" className="form-label">Chat Instructions</label>
					<textarea onChange={handleChange} className="form-control" id="create-bot-form-instructions" name="instructions" rows="3" placeholder="Enter chat instructions" value={formData?.instructions ??  ''} ></textarea>
					<div className="help-text text-start">This prompt will be sent to the bot as instructions for how to act in the chat.</div>
				</div>
				<div className="mb-3">
					<label htmlFor="memory_duration" className="form-label">Memory <RequiredStar/></label>
					<select onChange={handleChange} id="create-chat-form-model" required name="memory_duration" className="form-select" value={formData?.memory_duration}>
						{Object.entries(memory_durations).map(([key, value]) => (
							<option key={key} value={value}>{key}</option>
						))}
					</select>
				</div>
				<div className="mb-3">
					<label htmlFor="bot_id" className="form-label">Bot <RequiredStar/></label>
					<select onChange={handleChange} id="create-chat-form-model" required name="bot_id" className="form-select" value={formData?.bot_id}>
						<option value="">Select Bot</option>
						{Object.values(bots).map(bot => (
							<option key={bot.ID} value={bot.ID}>
								{bot.name}{bot.description ? " - "+bot.description : ""}</option>
						))}
					</select>
				</div>
				<div>
					<div className="btn border me-3" onClick={() => dispatch(goToChatPage(selectedChat))}>Cancel</div>
					<button type="submit" className="btn btn-primary">{selectedChat ? "Save" : "Create"} Chat</button>
				</div>
			</form>
		</div>
	);
}