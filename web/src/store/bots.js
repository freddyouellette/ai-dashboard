import { createSelector, createSlice } from '@reduxjs/toolkit'
import axios from 'axios';

const botsSlice = createSlice({
	name: 'bots',
	initialState: {
		bots: {},
		botsLoading: false,
		botsError: null,
		botsLoaded: false,
		botModels: {},
		botModelsLoading: false,
		botModelsError: null,
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
		setBotModels: (state, action) => {
			state.botModels = action.payload;
		},
		setBotModelsLoading: (state, action) => {
			state.botModelsLoading = action.payload;
		},
		setBotModelsLoaded: (state, action) => {
			state.botModelsLoaded = action.payload;
		},
		setBotModelsError: (state, action) => {
			state.botModelsError = action.payload;
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

export const getBotModels = () => async (dispatch, getState) => {
	if (getState().bots.botModelsLoaded) {
		return;
	}
	dispatch(botsSlice.actions.setBotModelsLoading(true));
	dispatch(botsSlice.actions.setBotModelsError(null));
	axios.get(process.env.REACT_APP_API_HOST+'/api/bots/models')
	.then(
		res => {
			let botModelsById = {};
			res.data.forEach(model => {
				botModelsById[model.id] = model;
			});
			dispatch(botsSlice.actions.setBotModelsLoading(false));
			dispatch(botsSlice.actions.setBotModelsLoaded(true));
			dispatch(botsSlice.actions.setBotModelsError(null));
			dispatch(botsSlice.actions.setBotModels(botModelsById));
		}, 
		error => {
			dispatch(botsSlice.actions.setBotModelsLoading(false));
			dispatch(botsSlice.actions.setBotModelsError(error));
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

export const selectBotModels = createSelector(
	state => state.bots.botModels,
	state => state.bots.botModelsLoading,
	state => state.bots.botModelsError, 
	(botModels, botModelsLoading, botModelsError) => ({ botModels, botModelsLoading, botModelsError })
);

export default botsSlice.reducer;