import { QueryClient, QueryClientProvider } from 'react-query';

import { ChakraProvider } from '@chakra-ui/react';
import { ApolloProvider } from '@apollo/client';
import { getApolloClient } from '../apollo';

const queryClient = new QueryClient();

const App = ({ Component, pageProps }) => {
  return (
    <ChakraProvider>
      <QueryClientProvider client={queryClient}>
        <ApolloProvider client={getApolloClient()}>
          <Component {...pageProps} />
        </ApolloProvider>
      </QueryClientProvider>
    </ChakraProvider>
  );
};

export default App;
