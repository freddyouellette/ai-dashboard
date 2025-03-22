// import Markdown from "react-markdown";
import remarkMath from "remark-math";
import rehypeKatex from "rehype-katex";
import "katex/dist/katex.min.css";
import {marked} from "marked";

export default function MessageText({ children }) {
	// ChatGPT returns \( and \), and sometimes \[ and \] instead of $ for LaTeX format
	// Replace \( and \) with $ for inline LaTeX, but ignore backticks
	children = children.replace(/\\+[()[\]](?=(?:[^`]*`[^`]*`)*[^`]*$)/g, '$');
	// Replace \[ and \] with $$ for block LaTeX, but ignore backticks
	children = children.replace(/\\+[[\]](?=(?:[^`]*`[^`]*`)*[^`]*$)/g, '$$');
	children = children.replace(/<think>/g, "<details class=\"text-muted\"><summary><i>Thinking...</i></summary><i>");
	children = children.replace(/<\/think>/g, "</i></details>");

	// let hasReplacedFirstThink = false;
	return <div dangerouslySetInnerHTML={{__html: marked.parse(children)}}/>
}