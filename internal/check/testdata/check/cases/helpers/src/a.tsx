import { cn } from "x";
export const A = () => <div className={cn("tg-row md:tg-gap(md)", f(g(1), (2)), `tg-col ${cond ? "tg-wrap" : "bad-one"}`)}>x</div>;
export const B = () => <div className={cn({ "bad-two": ok, "tg-fixed": !ok }, [`x-${n}`])}>x</div>;
export const C = () => <div className="tg-row" />;
