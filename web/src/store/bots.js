import { createSlice } from '@reduxjs/toolkit'

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
		}
	}
});

export const fetchBots = () => async dispatch => {
	fetch('http://localhost:8080/bots')
	.then(response => response.json())
	.then(
		bots => dispatch(botsSlice.actions.setBots(bots)), 
		error => console.error(error)
	);
}

export const { addOrUpdateBot, setBots } = botsSlice.actions;
export const selectBots = state => state.bots.bots;
export default botsSlice.reducer;