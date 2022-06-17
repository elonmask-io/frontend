import {Injectable} from '@angular/core';
import axios from "axios";
import {environment} from "../../environments/environment";
import {JwtHelperService} from "@auth0/angular-jwt";
@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {
  private jwtHelper : JwtHelperService
  constructor() {
    this.jwtHelper = new JwtHelperService()
  }

  async challenge(): Promise<string> {
    return "challenge"
  }

  async response(credential: Credential): Promise<any> {

  }

  async registerInitialize(username: string, name: string): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/register/initialize",
      JSON.stringify({username: username, name: name}),
      {
        headers: {
          "x-username": username
        }
      }
    )
  }

  async registerFinalize(username: string, token: any): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/register/finalize",
      token,
      {
        headers: {
          "x-username": username
        }
      }
    )
  }

  async loginInitialize(username: string): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/login/initialize",
      null,
      {
        headers: {
          "x-username": username
        }
      }
    )
  }

  async loginFinalize(username: string, token: any): Promise<string> {
    return axios.post(
      environment.routes.authenticationService + "/login/finalize",
      token,
      {
        headers: {
          "x-username": username
        }
      }
    )
  }

  public isAuthenticated(): boolean {
    const token = localStorage.getItem('token');
    if (token == null){
      return false
    }
    return !this.jwtHelper.isTokenExpired(token);
  }

  public login(token : string) : boolean {
    localStorage.setItem("token", token)
    return true
  }




}

