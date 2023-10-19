import { useDispatch, useSelector } from "react-redux";
import { selectBots } from "../store/bots";
import { Button, Container, ListGroup, ListGroupItem } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowRight, faPlus } from "@fortawesome/free-solid-svg-icons";
import { setBotSelected, goToBotEdit } from "../store/page";

export default function SidebarBotList() {
	const dispatch = useDispatch();
	const bots = useSelector(selectBots);
	
	return (
		<ListGroup className="list-group-flush">
			<ListGroupItem className="bg-dark border-bottom">
				<Container className="text-center">
					<Button onClick={() => {dispatch(goToBotEdit(null))}}>
						<FontAwesomeIcon icon={faPlus} className="me-2" />
						Add New Bot
					</Button>
				</Container>
			</ListGroupItem>
			{bots.map(bot => {
				return (
					<ListGroupItem 
						key={bot.ID} 
						className="bg-dark text-white border-bottom" 
						style={{ "cursor": "pointer" }}
						onClick={() => dispatch(setBotSelected(bot))}>
						<Container className="text-start">
							<div className="d-flex justify-content-between">
								<div className="flex-grow-1">
									<strong>{bot.name}</strong>
									<div>{bot.description}</div>
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