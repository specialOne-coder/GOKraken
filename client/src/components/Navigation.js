import React from "react";
import {
  Routes,
  Route,
  HashRouter as Router,
} from "react-router-dom";
import { HomePage,Archives,Pairs,UserArchive } from "../pages/index";
import {Navbar} from "./index";


const Navigation = () => {
  return (
    <Router>
      <Navbar />
      <Routes>
        <Route exact path="/" element={<HomePage />} />
        <Route exact path="/archive" element={<UserArchive />} />
        <Route exact path="/pairs" element={<Pairs />} />
        <Route exact path="/archives" element={<Archives />} />
      </Routes>
    </Router>
  );
};


export default Navigation;
