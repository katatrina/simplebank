import type { AuthState } from "@/types/auth_state";
import type { User } from "@/types/user";
import { reactive, readonly } from "vue";

const state = reactive<AuthState>({
  user: null,
  access_token: null,
  refresh_token: null,
})

function setUser(user: User, access_token: string, refresh_token: string) {
  state.user = user;
  state.access_token = access_token;
  state.refresh_token = refresh_token;
}

function clearUser() {
  state.user = null;
  state.access_token = null;
  state.refresh_token = null;
}

export default {
  state: readonly(state),
  setUser,
  clearUser,
}
