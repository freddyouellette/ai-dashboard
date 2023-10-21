import { createSlice } from '@reduxjs/toolkit'
import axios from 'axios';

const chatsSlice = createSlice({
	name: 'chats',
	initialState: {
		chats: [],
		selectedChat: null,
	},
	reducers: {
		addChat: (state, action) => {
			state.chats.push(action.payload)
		},
		setChats: (state, action) => {
			state.chats = action.payload;
		},
		setSelectedChat: (state, action) => {
			state.selectedChat = action.payload
		},
	}
});

export const fetchChats = () => async dispatch => {
	fetch('http://localhost:8080/chats')
	.then(response => response.json())
	.then(
		chats => dispatch(chatsSlice.actions.setChats(chats)), 
		error => console.error(error)
	);
}

// thunk
export const addChat = (botId) => async dispatch => {
	let newChatData = {
		name: "New Chat",
		bot_id: botId,
	}
	axios.post('http://localhost:8080/chats', newChatData)
	.then(response => {
		console.log(response);
		dispatch(chatsSlice.actions.addChat(response.data));
	}).catch(error => {
		console.error(error);
	})
}

export const { setSelectedChat, setChats } = chatsSlice.actions;
export const selectChats = state => state.chats.chats;
export const selectSelectedChat = state => state.chats.selectedChat;
export default chatsSlice.reducer;