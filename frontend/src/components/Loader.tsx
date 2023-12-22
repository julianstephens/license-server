import { ReactSVG } from "react-svg";

export const Loader = () => {
  return (
    <div className="col centered my-10">
      <ReactSVG src="/spinner.svg" />
      <h4>Loading...</h4>
    </div>
  );
};
