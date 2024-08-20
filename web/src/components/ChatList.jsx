import { useDispatch, useSelector } from "react-redux";
import { Container, ListGroup } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronRight } from "@fortawesome/free-solid-svg-icons";
import { goToChatPage } from "../store/page";
import { getChats, selectChats } from "../store/chats";
import { useEffect } from "react";
import { getBots, selectBots } from "../store/bots";
import moment from "moment";

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
	
	// put chats in order of creation descending
	let chatList = Object.values(chats);
	chatList.sort((a, b) => moment(b.last_message_at || b.CreatedAt) - moment(a.last_message_at || a.CreatedAt));
	
	return (
		<ListGroup className="list-group-flush">
			{chatList.map(chat => {
				let bot = bots[chat.bot_id];
				return (
					<div 
						key={chat.ID} 
						className="border-bottom" 
						style={{ "cursor": "pointer" }}
						onClick={() => dispatch(goToChatPage(chat))}>
						<div className="text-start mx-3 my-1">
							<div className="d-flex justify-content-between">
								<div className="flex-grow-1">
									<div><strong>{chat.name || chat.ID}</strong></div>
									<div>{bot?.name ?? <i>unknown</i>}</div>
									<div className="text-muted"><i><small>{bot?.model}</small></i></div>
								</div>
								<div className="d-flex align-items-center">
									<FontAwesomeIcon 
										icon={faChevronRight} 
										className="ms-2 cursor-pointer" 
										style={{"cursor": "pointer"}} 
										/>
								</div>
							</div>
						</div>
					</div>
				);
			})}
		</ListGroup>
	);
}