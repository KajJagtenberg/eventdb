import {
  ApolloClient,
  HttpLink,
  InMemoryCache,
  NormalizedCacheObject,
} from '@apollo/client';
import { backend } from './vars/backend';

const createApolloClient = (): ApolloClient<NormalizedCacheObject> => {
  return new ApolloClient({
    ssrMode: typeof window === 'undefined',
    link: new HttpLink({
      uri: `${backend}/query`,
    }),
    cache: new InMemoryCache(),
  });
};

let client: ApolloClient<NormalizedCacheObject>;

export const getApolloClient = (): ApolloClient<NormalizedCacheObject> => {
  if (client == undefined) {
    client = createApolloClient();
  }

  return client;
};
