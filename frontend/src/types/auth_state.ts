import type { User } from "@/types/user";

export interface AuthState {
  user: User | null,
  access_token: string | null,
  refresh_token: string | null,
}
