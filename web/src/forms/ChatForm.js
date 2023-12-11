import { useDispatch, useSelector } from "react-redux";
import { getBots, selectBots } from "../store/bots";
import RequiredStar from './RequiredStar';
import { createChat } from "../store/chats";
import { useEffect, useState } from "react";

export default function ChatForm() {
	const dispatch = useDispatch();
	const { bots, botsLoading, botsError } = useSelector(selectBots);
	const [formData, setFormData] = useState({
		botId: null,
	});
	console.log(formData);
	useEffect(() => {
		dispatch(getBots());
	}, [dispatch]);
	
	if (botsLoading) return <div>Loading...</div>
	if (botsError) return <div>Error loading bots...</div>
	
	const handleSubmit = async (event) => {
		event.preventDefault();
		
		let createChatData = Object.assign({}, formData);
		createChatData.botId = parseInt(createChatData.botId);
		dispatch(createChat(createChatData));
	}
	
	const handleChange = (event) => {
		setFormData({
			...formData,
			[event.target.name]: event.target.value
		});
	}
	
	return (
		<div className="container">
			<h1 className="text-center">New Chat</h1>
			<form onSubmit={handleSubmit}>
				<div className="mb-3">
					<label htmlFor="botId" className="form-label">Bot <RequiredStar/></label>
					<select onChange={handleChange} id="create-chat-form-model" required name="botId" className="form-control">
						<option value="">Select Bot</option>
						{Object.values(bots).map(bot => (
							<option key={bot.ID} value={bot.ID}>{bot.name}</option>
						))}
					</select>
				</div>
				<button type="submit" className="btn btn-primary">Create Chat</button>
			</form>
		</div>
	);
}