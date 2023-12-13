import { createSelector, createSlice } from '@reduxjs/toolkit'
import axios from 'axios';

const botsSlice = createSlice({
	name: 'bots',
	initialState: {
		bots: {},
		botsLoading: false,
		botsError: null,
		botsLoaded: false,
	},
	reducers: {
		setBots: (state, action) => {
			state.bots = action.payload;
		},
		setBotsLoading: (state, action) => {
			state.botsLoading = action.payload;
		},
		setBotsLoaded: (state, action) => {
			state.botsLoaded = action.payload;
		},
		setBotsError: (state, action) => {
			state.botsError = action.payload;
		},
	}
});

export const getBots = () => async (dispatch, getState) => {
	if (getState().bots.botsLoaded) {
		return;
	}
	dispatch(botsSlice.actions.setBotsLoading(true));
	dispatch(botsSlice.actions.setBotsError(null));
	axios.get(process.env.REACT_APP_API_HOST+'/api/bots')
	.then(
		res => {
			let botsById = {};
			res.data.forEach(bot => {
				botsById[bot.ID] = bot;
			});
			dispatch(botsSlice.actions.setBotsLoading(false));
			dispatch(botsSlice.actions.setBotsError(null));
			dispatch(botsSlice.actions.setBots(botsById));
			dispatch(botsSlice.actions.setBotsLoaded(true));
		}, 
		error => {
			dispatch(botsSlice.actions.setBotsLoading(false));
			dispatch(botsSlice.actions.setBotsLoaded(false));
			dispatch(botsSlice.actions.setBotsError(error));
			console.error(error);
		},
	);
}

export const addOrUpdateBot = (bot) => async (dispatch, getState) => {
	dispatch(botsSlice.actions.setBotsLoading(true));
	dispatch(botsSlice.actions.setBotsError(null));
	return axios[bot.ID ? "put" : "post"](process.env.REACT_APP_API_HOST+"/api/bots", bot)
	.then(
		res => {
			dispatch(botsSlice.actions.setBotsLoading(false));
			dispatch(botsSlice.actions.setBotsError(null));
			let oldBots = getState().bots.bots;
			dispatch(botsSlice.actions.setBots({ ...oldBots, [res.data.ID]: res.data }));
			return res.data;
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