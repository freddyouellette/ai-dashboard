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
	});
	useEffect(() => {
		dispatch(getBots());
	}, [dispatch]);
	
	if (botsLoading) return <div>Loading...</div>
	if (botsError) return <div>Error loading bots...</div>
	
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
					<label htmlFor="bot_id" className="form-label">Bot <RequiredStar/></label>
					<select onChange={handleChange} id="create-chat-form-model" required name="bot_id" className="form-control" value={formData?.bot_id}>
						<option value="">Select Bot</option>
						{Object.values(bots).map(bot => (
							<option key={bot.ID} value={bot.ID}>
								{bot.name}{bot.description ? " - "+bot.description : ""}</option>
						))}
					</select>
				</div>
				<button type="submit" className="btn btn-primary">{selectedChat ? "Save" : "Create"} Chat</button>
			</form>
		</div>
	);
}