import { useDispatch, useSelector } from "react-redux";
import { getBots, selectBots } from "../store/bots";
import { Container, ListGroup, ListGroupItem } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowRight } from "@fortawesome/free-solid-svg-icons";
import { goToBotEditPage } from "../store/page";
import moment from "moment";
import { useEffect } from "react";

export default function BotList() {
	const dispatch = useDispatch();
	const { bots, botsLoading, botsError } = useSelector(selectBots);
	
	useEffect(() => {
		dispatch(getBots());
	}, [dispatch]);
	
	if (botsLoading) return <div>Loading...</div>;
	if (botsError) return <div className="text-danger">Error loading bots...</div>;
	
	let botsList = Object.values(bots);
	botsList.sort((a, b) => moment(b.last_message_at || b.CreatedAt) - moment(a.last_message_at || a.CreatedAt));
	
	return (
		<ListGroup className="list-group-flush">
			<ListGroupItem className="border-bottom">
				<div className="d-flex justify-content-center">
					<div className="btn btn-light border" onClick={() => dispatch(goToBotEditPage())}>Create New Bot</div>
				</div>
			</ListGroupItem>
			{botsList.map(bot => {
				return (
					<ListGroupItem 
						key={bot.ID} 
						className="border-bottom" 
						style={{ "cursor": "pointer" }}
						onClick={() => dispatch(goToBotEditPage(bot))}>
						<Container className="text-start">
							<div className="d-flex justify-content-between">
								<div className="flex-grow-1">
									<div><strong>{bot.name}</strong></div>
									<div><div>{bot.description}</div></div>
									<div className="text-muted"><i><small>{bot.model}</small></i></div>
								</div>
								<div className="d-flex align-items-center">
									<FontAwesomeIcon 
										icon={faArrowRight} 
										className="ms-2 cursor-pointer" 
										style={{"cursor": "pointer"}} 
										/>
								</div>
							</div>
						</Container>
					</ListGroupItem>
				);
			})}
		</ListGroup>
	);
}