import { useDispatch, useSelector } from "react-redux";
import { Container, ListGroup, ListGroupItem } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronRight, faPlus } from "@fortawesome/free-solid-svg-icons";
import { goToChatPage } from "../store/page";
import { getChats, selectChats } from "../store/chats";
import { TODO } from "../util/todo";
import { useEffect } from "react";
import { getBots, selectBots } from "../store/bots";

export default function ChatList() {
	const dispatch = useDispatch();
	useEffect(() => {
		dispatch(getChats());
		dispatch(getBots());
	}, [dispatch]);
	let { chats, chatsLoading, chatsError } = useSelector(selectChats);
	let { bots, botsLoading, botsError } = useSelector(selectBots);
	
	if (chatsLoading || botsLoading) return <div>Loading...</div>;
	if (chatsError) return <div>Error loading chats...</div>;
	if (botsError) return <div>Error loading bots...</div>;
	
	return (
		<ListGroup className="list-group-flush">
			{Object.values(chats).map(chat => {
				let bot = bots[chat.bot_id];
				return (
					<div 
						key={chat.ID} 
						className="border-bottom bg-light" 
						style={{ "cursor": "pointer" }}
						onClick={() => dispatch(goToChatPage(chat))}>
						<Container className="text-start">
							<div className="d-flex justify-content-between">
								<div className="flex-grow-1">
									<div><strong>{bot?.name ?? <i>unknown</i>}</strong></div>
									<div>{chat.ID}</div>
								</div>
								<div className="d-flex align-items-center">
									<FontAwesomeIcon 
										icon={faChevronRight} 
										className="ms-2 cursor-pointer" 
										style={{"cursor": "pointer"}} 
										/>
								</div>
							</div>
						</Container>
					</div>
				);
			})}
		</ListGroup>
	);
}