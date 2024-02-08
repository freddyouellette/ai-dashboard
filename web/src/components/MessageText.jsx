import Markdown from "react-markdown";
import remarkMath from "remark-math";
import rehypeKatex from "rehype-katex";
import "katex/dist/katex.min.css";

export default function MessageText({ children }) {
	// ChatGPT returns \( and \), and sometimes \[ and \] instead of $ for LaTeX format
	// Replace \( and \) with $ for inline LaTeX, but ignore backticks
	children = children.replace(/\\+[()[\]](?=(?:[^`]*`[^`]*`)*[^`]*$)/g, '$');
	// Replace \[ and \] with $$ for block LaTeX, but ignore backticks
	children = children.replace(/\\+[[\]](?=(?:[^`]*`[^`]*`)*[^`]*$)/g, '$$');
	
	return (
		<Markdown remarkPlugins={[remarkMath]} rehypePlugins={[rehypeKatex]} className="text-break">
			{children}
		</Markdown>
	)
}