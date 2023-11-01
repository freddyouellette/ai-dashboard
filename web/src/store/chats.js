import { createSlice } from '@reduxjs/toolkit'
import axios from 'axios';
import { goToBotChat } from './page';

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
	}
});

// thunk
export const fetchChats = () => async dispatch => {
	fetch('http://localhost:8080/api/chats')
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
	axios.post('http://localhost:8080/api/chats', newChatData)
	.then(response => {
		console.log(response);
		dispatch(chatsSlice.actions.addChat(response.data));
		dispatch(goToBotChat(response.data))
	}).catch(error => {
		console.error(error);
	})
}

export const { setChats } = chatsSlice.actions;
export const selectChats = state => state.chats.chats;
export default chatsSlice.reducer;