import { ChakraProvider } from "@chakra-ui/react";
import { QueryClientProvider, QueryClient } from "react-query";

const queryClient = new QueryClient();

const App = ({ Component, pageProps }) => {
  return (
    <ChakraProvider>
      <QueryClientProvider client={queryClient}>
        <Component {...pageProps} />
      </QueryClientProvider>
    </ChakraProvider>
  );
};

export default App;
