import type { AuthState } from "@/types/auth_state";
import type { User } from "@/types/user";
import { reactive, readonly } from "vue";

const authState = reactive<AuthState>({
  user: null,
  access_token: null,
  refresh_token: null,
})

function setUser(user: User, access_token: string, refresh_token: string) {
  authState.user = user;
  authState.access_token = access_token;
  authState.refresh_token = refresh_token;
}

function clearUser() {
  authState.user = null;
  authState.access_token = null;
  authState.refresh_token = null;
}

export default {
  authState: readonly(authState),
  setUser,
  clearUser,
}
