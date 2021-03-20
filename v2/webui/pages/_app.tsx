import { ApolloProvider, ApolloClient, InMemoryCache } from '@apollo/client';
import { ChakraProvider } from '@chakra-ui/react';
import { theme } from '../theme';

const client = new ApolloClient({
  uri: 'http://127.0.0.1:16543/graphql',
  cache: new InMemoryCache(),
});

const App = ({ Component, pageProps }) => {
  return (
    <ChakraProvider theme={theme}>
      <ApolloProvider client={client}>
        <Component {...pageProps} />
      </ApolloProvider>
    </ChakraProvider>
  );
};

export default App;
