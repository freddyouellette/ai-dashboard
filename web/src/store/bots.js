import { createSlice } from '@reduxjs/toolkit'

const botsSlice = createSlice({
	name: 'bots',
	initialState: {
		bots: []
	},
	reducers: {
		upsertBot: (state, action) => {
			const bot = action.payload;
			const botIndex = state.bots.findIndex(b => b.ID === bot.ID);
			if(botIndex > -1) {
				// Update bot
				state.bots[botIndex] = bot;
			} else {
				// Create bot
				state.bots.push(bot);
			}
		}
	}
});

export const { upsertBot } = botsSlice.actions;
export const { selectBots } = state => state.bots;
export default botsSlice;