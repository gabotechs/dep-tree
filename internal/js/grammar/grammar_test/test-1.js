import React from "react"
import { CustomApolloClient } from "../services/apollo";
import config from "../config";
import createTheme, { PADDING_EX } from "../styles/theme";
import { createBrowserHistory } from "history";
import { ROUTING_LOGIN } from "../constants/routing";
import { LogProvider } from "../components/dialogs/SnackBar";
import { ApolloProvider } from "@apollo/client";
import { Router as RouterProvider } from "react-router";
import { AppContextProvider } from "./contexts/AppContext";
import { MuiThemeProvider } from "@material-ui/core";
import { SlicingContextProvider } from "../views/SlicerView/contexts/SlicingContext";
