import { createSlice } from '@reduxjs/toolkit'

const chatsSlice = createSlice({
	name: 'chats',
	initialState: {
		chats: []
	},
	reducers: {
		setChats: (state, action) => {
			state.chats = action.payload;
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

export const { addOrUpdateBot, setChats } = chatsSlice.actions;
export const selectChats = state => state.chats.chats;
export default chatsSlice.reducer;