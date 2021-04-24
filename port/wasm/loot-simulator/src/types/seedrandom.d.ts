// Adapted from https://github.com/davidbau/seedrandom/pull/70/files#diff-093ad82a25aee498b11febf1cdcb6546e4d223ffcb49ed69cc275ac27ce0ccce

declare module "seedrandom" {
  const seedrandom: Seedrandom;
  export default seedrandom;

  export interface Seedrandom {
    (seed?: any, options?: boolean): prng.Prng;
    <TOptions extends option.Seedrandom>(
      seed: any,
      options: TOptions
    ): TOptions extends option.Pass<infer TReturn>
      ? TReturn
      : TOptions extends option.Global
      ? string
      : TOptions extends option.State<state.Arc4>
      ? prng.WithState<state.Arc4>
      : prng.Prng;
    <TOptions extends option.Seedrandom | boolean, TReturn>(
      seed: any,
      options: TOptions,
      callback: Callback<TReturn>
    ): TOptions extends option.Pass<infer TReturn> ? TReturn : TReturn;

    alea: AlterMethod<state.Alea>;
    xor128: AlterMethod<state.Xor128>;
    xorwow: AlterMethod<state.Xorwow>;
    xorshift7: AlterMethod<state.Xorshift7>;
    xor4096: AlterMethod<state.Xor4096>;
    tychei: AlterMethod<state.Tychei>;
  }

  export type Callback<TReturn> = (prng: prng.Union<state.Arc4>, seed: string) => TReturn;

  export interface AlterMethod<TState> {
    (seed?: any): prng.Prng;
    (seed: any, opts?: option.State<TState>): prng.WithState<TState>;
  }

  export namespace prng {
    interface Prng {
      (): number;
      int32(): number;
      quick(): number;
      double(): number;
    }
    interface WithState<TState> extends Prng {
      state(): TState;
    }
    type Union<TState> = Prng | WithState<TState>;
  }

  export namespace state {
    type State<Props extends string> = { [P in Props]: number };
    type StateWithArray<Props extends string, ArrayProps extends string> = State<Props> &
      { [P in ArrayProps]: number[] };

    type Arc4 = StateWithArray<"i" | "j", "S">;
    type Alea = State<"c" | "s0" | "s1" | "s2">;
    type Xor128 = State<"x" | "y" | "z" | "w">;
    type Xorwow = State<"x" | "y" | "z" | "w" | "v" | "d">;
    type Xorshift7 = StateWithArray<"i", "x">;
    type Xor4096 = StateWithArray<"i" | "w", "X">;
    type Tychei = State<"a" | "b" | "c" | "d">;
  }

  export namespace option {
    interface Entropy {
      entropy: true;
    }
    interface State<TState> {
      state: TState | true;
    }
    interface Global {
      global: true;
    }
    interface Pass<TReturn> {
      pass: Callback<TReturn>;
    }

    type Math = Partial<Entropy & State<state.Arc4>>;
    type Seedrandom = Partial<Math & Global & Pass<any>>;
  }
}
