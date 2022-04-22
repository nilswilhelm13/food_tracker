import Dashboard from "../components/Dashboard/Dashboard";
import Logout from "../components/Auth/Logout";
import Scan from "@/components/Scan";
import Auth from "@/components/Auth/Auth";
import AddWrapper from "@/components/NewFood/AddWrapper";


export const routes = [
	{path: "", component: Dashboard},
	{path: "/addfood", component: AddWrapper},
	{path: "/login", component: Auth},
	{path: "/logout", component: Logout},
	{path: "/scan", component: Scan},
];