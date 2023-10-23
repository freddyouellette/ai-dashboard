import { createSlice } from '@reduxjs/toolkit'
import axios from 'axios';
import { goToSidebarBotList } from './page';

const botsSlice = createSlice({
	name: 'bots',
	initialState: {
		bots: []
	},
	reducers: {
		addOrUpdateBot: (state, action) => {
			const bot = action.payload;
			const botIndex = state.bots.findIndex(b => b.ID === bot.ID);
			if(botIndex > -1) {
				// Update bot
				state.bots[botIndex] = bot;
			} else {
				// Create bot
				state.bots.push(bot);
			}
		},
		setBots: (state, action) => {
			state.bots = action.payload;
		},
		removeBot: (state, action) => {
			const bot = action.payload;
			const botIndex = state.bots.findIndex(b => b.ID === bot.ID);
			if(botIndex > -1) {
				state.bots.splice(botIndex, 1);
			}
		}
	}
});

export const fetchBots = () => async dispatch => {
	axios.get('http://localhost:8080/bots')
	.then(
		res => dispatch(botsSlice.actions.setBots(res.data)), 
		error => console.error(error)
	);
}

export const deleteBot = bot => async (dispatch, getState) => {
	axios.delete(`http://localhost:8080/bots/${bot.ID}`)
	.then(
		() => {
			dispatch(botsSlice.actions.removeBot(bot))
			const state = getState();
			if (state.page.sidebarBotSelected?.ID === bot.ID) {
				dispatch(goToSidebarBotList());
			}
		},
		error => console.error(error)
	);
}

export const { addOrUpdateBot } = botsSlice.actions;
export const selectBots = state => state.bots.bots;
export default botsSlice.reducer;