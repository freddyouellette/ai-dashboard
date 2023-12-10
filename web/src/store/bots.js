import { createSelector, createSlice } from '@reduxjs/toolkit'
import axios from 'axios';

const botsSlice = createSlice({
	name: 'bots',
	initialState: {
		bots: {},
		botsLoading: false,
		botsError: null,
	},
	reducers: {
		setBots: (state, action) => {
			state.bots = action.payload;
		},
		setBotsLoading: (state, action) => {
			state.botsLoading = action.payload;
		},
		setBotsError: (state, action) => {
			state.botsError = action.payload;
		},
	}
});

export const getBots = () => async dispatch => {
	dispatch(botsSlice.actions.setBotsLoading(true));
	dispatch(botsSlice.actions.setBotsError(null));
	axios.get('http://localhost:8080/api/bots')
	.then(
		res => {
			let botsById = {};
			res.data.forEach(bot => {
				botsById[bot.ID] = bot;
			});
			dispatch(botsSlice.actions.setBotsLoading(false));
			dispatch(botsSlice.actions.setBotsError(null));
			dispatch(botsSlice.actions.setBots(botsById));
		}, 
		error => {
			dispatch(botsSlice.actions.setBotsLoading(false));
			dispatch(botsSlice.actions.setBotsError(error));
			console.error(error);
		},
	);
}

export const selectBots = createSelector(
	state => state.bots.bots,
	state => state.bots.botsLoading,
	state => state.bots.botsError, 
	(bots, botsLoading, botsError) => ({ bots, botsLoading, botsError })
);

export default botsSlice.reducer;