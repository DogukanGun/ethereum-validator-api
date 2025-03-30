"use client";
import { useState } from "react";

interface SyncData {
  validators: string[] | null;
  sync_info: {
    sync_period: number;
    committee_size: number;
  };
}

interface BlockData {
  status: string;
  reward: number;
  block_info: {
    proposer_payment: number;
    is_mev_boost: boolean;
  };
}

interface ApiError {
  error: string;
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'https://sf-api.dogukangun.de/';

const getApiUrl = () => {
  if (typeof window !== 'undefined') {
    return 'https://sf-api.dogukangun.de/';
  }
  return API_URL;
};

export default function Home() {
  const [slot, setSlot] = useState("");
  const [syncData, setSyncData] = useState<SyncData | null>(null);
  const [blockData, setBlockData] = useState<BlockData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = async () => {
    if (!slot) {
      setError("Please enter a slot number");
      return;
    }
    
    setLoading(true);
    setError(null);
    setSyncData(null);
    setBlockData(null);

    try {
      const apiUrl = getApiUrl();
      
      // Attempt to fetch data from both endpoints
      const syncRes = await fetch(`${apiUrl}/syncduties/${slot}`, {
        headers: { 'Accept': 'application/json' },
      });
      
      const blockRes = await fetch(`${apiUrl}/blockreward/${slot}`, {
        headers: { 'Accept': 'application/json' },
      });

      // Process sync committee data
      let syncJson = null;
      if (syncRes.ok) {
        syncJson = await syncRes.json() as SyncData;
      } else {
        const syncErrorData = await syncRes.json() as ApiError;
        console.warn('Sync data fetch failed:', syncErrorData);
        
        // Create a placeholder response to show something in the UI
        syncJson = {
          validators: [],
          sync_info: {
            sync_period: 0,
            committee_size: 0,
          }
        };
      }

      // Process block reward data
      let blockJson = null;
      if (blockRes.ok) {
        blockJson = await blockRes.json() as BlockData;
      } else {
        const blockErrorData = await blockRes.json() as ApiError;
        console.warn('Block data fetch failed:', blockErrorData);
        
        // Create a placeholder response to show something in the UI
        blockJson = {
          status: "unknown",
          reward: 0,
          block_info: {
            proposer_payment: 0,
            is_mev_boost: false,
          }
        };
        
        if (blockErrorData && blockErrorData.error) {
          setError(`Block data: ${blockErrorData.error}`);
        }
      }

      setSyncData(syncJson);
      setBlockData(blockJson);
    } catch (error) {
      console.error('Error fetching data:', error);
      setError(error instanceof Error ? error.message : 'Network error - please check API connection');
      
      // Create placeholder data so UI doesn't display nulls
      setSyncData({
        validators: [],
        sync_info: {
          sync_period: 0, 
          committee_size: 0
        }
      });
      
      setBlockData({
        status: "error",
        reward: 0,
        block_info: {
          proposer_payment: 0,
          is_mev_boost: false
        }
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 text-white p-8">
      <main className="max-w-6xl mx-auto">
        <h1 className="text-4xl font-bold text-center mb-4 text-white drop-shadow-lg">
          Ethereum Validator Explorer
        </h1>
        
        <p className="text-center text-indigo-300/70 mb-8 max-w-2xl mx-auto">
          Explore Ethereum validator performance, sync committee duties, and block rewards across the Beacon Chain.
        </p>
        
        {/* Stats summary - shows when data is loaded */}
        {(syncData || blockData) && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
            <div className="bg-slate-800/40 rounded-lg p-4 border border-indigo-300/10 backdrop-blur-sm">
              <div className="text-sm text-indigo-300/80 mb-1">Slot</div>
              <div className="text-2xl font-bold text-white">{slot}</div>
              <div className="text-xs text-indigo-300/50 mt-1">
                {new Date(Number(slot) * 12 * 1000).toLocaleString()}
              </div>
            </div>
            
            <div className="bg-slate-800/40 rounded-lg p-4 border border-indigo-300/10 backdrop-blur-sm">
              <div className="text-sm text-indigo-300/80 mb-1">Sync Committee</div>
              <div className="text-2xl font-bold text-white">
                {syncData?.validators?.length || 0} Validators
              </div>
              <div className="text-xs text-indigo-300/50 mt-1">
                Period: {syncData?.sync_info?.sync_period || 'Unknown'}
              </div>
            </div>
            
            <div className="bg-slate-800/40 rounded-lg p-4 border border-indigo-300/10 backdrop-blur-sm">
              <div className="text-sm text-indigo-300/80 mb-1">Block Status</div>
              <div className={`text-2xl font-bold ${
                blockData?.status === 'mev' ? 'text-emerald-300' : 
                blockData?.status === 'vanilla' ? 'text-sky-300' : 'text-white'
              }`}>
                {blockData?.status === 'mev' ? 'MEV Block' : 
                 blockData?.status === 'vanilla' ? 'Vanilla Block' : 'Unknown'}
              </div>
              <div className="text-xs text-indigo-300/50 mt-1">
                Reward: {blockData?.reward?.toLocaleString() || 0} GWEI
              </div>
            </div>
          </div>
        )}
        
        <div className="flex flex-col items-center gap-4 mb-12">
          <div className="flex gap-4">
            <input
              type="number"
              value={slot}
              onChange={(e) => setSlot(e.target.value)}
              placeholder="Enter slot number"
              className="px-4 py-2 rounded-lg bg-slate-800 border border-indigo-300/30 focus:outline-none focus:ring-2 focus:ring-indigo-400/50"
            />
            <button
              onClick={fetchData}
              disabled={loading}
              className="px-6 py-2 bg-indigo-400 rounded-lg hover:bg-indigo-500 transition-colors disabled:opacity-50 flex items-center"
            >
              {loading ? (
                <>
                  <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Loading...
                </>
              ) : "Fetch Data"}
            </button>
          </div>
          <div className="text-sm text-indigo-300/70 max-w-lg text-center">
            Try these slots: <button onClick={() => setSlot("4700000")} className="text-indigo-200 hover:underline">4700000</button>, <button onClick={() => setSlot("4800000")} className="text-indigo-200 hover:underline">4800000</button>, or <button onClick={() => setSlot("4900000")} className="text-indigo-200 hover:underline">4900000</button>. 
            <p className="mt-2">The Ethereum Beacon Chain uses slots for block production (one every 12 seconds).</p>
          </div>
          {error && (
            <div className="text-rose-300 text-sm bg-rose-500/10 px-4 py-2 rounded-lg">
              {error}
            </div>
          )}
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          {/* Sync Committee Data */}
          <div className="bg-slate-800/50 rounded-xl p-6 border border-indigo-300/20 backdrop-blur-sm shadow-xl relative min-h-[320px]">
            <h2 className="text-2xl font-semibold mb-4 text-indigo-200">Sync Committee</h2>
            
            {loading ? (
              <div className="absolute inset-0 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm rounded-xl">
                <div className="text-indigo-200 flex flex-col items-center">
                  <svg className="animate-spin h-8 w-8 mb-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <span>Loading sync data...</span>
                </div>
              </div>
            ) : syncData && (
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-slate-300">Sync Period</span>
                  <span className="text-indigo-200">{syncData.sync_info.sync_period}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-slate-300">Committee Size</span>
                  <span className="text-indigo-200">{syncData.sync_info.committee_size}</span>
                </div>
                <div>
                  <h3 className="text-sm text-slate-300 mb-2 flex justify-between items-center">
                    <span>Validators</span>
                    {syncData.validators && (
                      <span className="bg-indigo-500/30 text-indigo-200 text-xs px-2 py-1 rounded">
                        {syncData.validators.length} total
                      </span>
                    )}
                  </h3>
                  <div className="max-h-[200px] overflow-y-auto bg-slate-900/50 rounded-lg p-4 scrollbar-thin scrollbar-thumb-indigo-500/30 scrollbar-track-transparent">
                    {syncData.validators && syncData.validators.length > 0 ? (
                      <div className="grid grid-cols-1 gap-2">
                        {syncData.validators.map((validator, i) => {
                          // Format validator key for better display
                          // Real validator keys are BLS public keys that are much longer
                          const shortKey = validator.length > 20 
                            ? `${validator.substring(0, 10)}...${validator.substring(validator.length - 6)}`
                            : validator;
                            
                          return (
                            <div 
                              key={i} 
                              className="text-xs flex items-center bg-slate-800/50 border border-indigo-500/10 p-2 rounded-lg hover:bg-slate-800/80 transition-colors group"
                              title={validator} // Show full key on hover
                            >
                              <div className="w-5 h-5 flex items-center justify-center bg-indigo-500/20 text-indigo-300 rounded-full mr-2 text-xs">
                                {i + 1}
                              </div>
                              <div className="flex flex-col flex-grow">
                                <span className="font-mono text-indigo-200 truncate">
                                  {shortKey}
                                </span>
                                <span className="text-xs text-slate-400">
                                  Validator #{i + 1}
                                </span>
                              </div>
                              <div className="opacity-0 group-hover:opacity-100 transition-opacity flex">
                                <a 
                                  href={`https://beaconcha.in/validator/${validator.replace('0x', '')}`} 
                                  target="_blank" 
                                  rel="noopener noreferrer"
                                  className="text-indigo-400 hover:text-indigo-300 p-1"
                                  title="View validator details on Beaconcha.in"
                                >
                                  <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                                  </svg>
                                </a>
                                <button 
                                  onClick={() => navigator.clipboard.writeText(validator)}
                                  className="text-indigo-400 hover:text-indigo-300 p-1 ml-1"
                                  title="Copy validator public key"
                                >
                                  <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
                                  </svg>
                                </button>
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    ) : (
                      <div className="text-slate-400 text-sm py-4 text-center">
                        <div className="flex flex-col items-center">
                          <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-slate-500 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                          </svg>
                          <span>No validators found for this slot</span>
                          <span className="text-xs text-indigo-400 mt-1">Try a different slot number</span>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Block Reward Data */}
          <div className="bg-slate-800/50 rounded-xl p-6 border border-indigo-300/20 backdrop-blur-sm shadow-xl relative min-h-[320px]">
            <h2 className="text-2xl font-semibold mb-4 text-indigo-200">Block Reward</h2>
            
            {loading ? (
              <div className="absolute inset-0 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm rounded-xl">
                <div className="text-indigo-200 flex flex-col items-center">
                  <svg className="animate-spin h-8 w-8 mb-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <span>Loading block data...</span>
                </div>
              </div>
            ) : blockData && (
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-slate-300">Status</span>
                  <span className={`px-3 py-1 rounded-full text-sm ${
                    blockData.status === 'mev' ? 'bg-emerald-400/20 text-emerald-200' : 
                    blockData.status === 'vanilla' ? 'bg-sky-400/20 text-sky-200' :
                    'bg-gray-400/20 text-gray-200'
                  }`}>
                    {blockData.status}
                  </span>
                </div>
                
                <div className="flex justify-between items-center">
                  <span className="text-slate-300">Reward (GWEI)</span>
                  <div className="flex items-center">
                    <div className="w-3 h-3 rounded-full bg-indigo-400 mr-2"></div>
                    <span className="text-indigo-200 font-medium">{blockData.reward.toLocaleString() || '0'}</span>
                  </div>
                </div>
                
                <div className="flex justify-between items-center">
                  <span className="text-slate-300">MEV-Boost</span>
                  <div className="flex items-center">
                    {blockData.block_info.is_mev_boost ? (
                      <>
                        <div className="w-3 h-3 rounded-full bg-emerald-400 mr-2"></div>
                        <span className="text-emerald-200">Yes</span>
                      </>
                    ) : (
                      <>
                        <div className="w-3 h-3 rounded-full bg-rose-400 mr-2"></div>
                        <span className="text-rose-200">No</span>
                      </>
                    )}
                  </div>
                </div>
                
                <div className="flex justify-between items-center">
                  <span className="text-slate-300">Proposer Payment</span>
                  <div className="flex items-center">
                    <div className="w-3 h-3 rounded-full bg-indigo-400 mr-2"></div>
                    <span className="text-indigo-200 font-medium">{blockData.block_info.proposer_payment.toLocaleString() || '0'}</span>
                  </div>
                </div>
                
                {/* Block explorer links */}
                <div className="flex justify-end pt-1">
                  <a 
                    href={`https://beaconcha.in/slot/${parseInt(slot)}`} 
                    target="_blank" 
                    rel="noopener noreferrer"
                    className="text-indigo-400 hover:text-indigo-300 text-xs flex items-center"
                  >
                    <span>View on Beaconcha.in</span>
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                    </svg>
                  </a>
                  <a 
                    href={`https://etherscan.io/block/${parseInt(slot)}`} 
                    target="_blank" 
                    rel="noopener noreferrer"
                    className="text-indigo-400 hover:text-indigo-300 text-xs flex items-center ml-4"
                  >
                    <span>View on Etherscan</span>
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                    </svg>
                  </a>
                </div>
                
                {/* Reward visualization */}
                <div className="pt-4">
                  <div className="h-10 bg-slate-900/50 rounded-lg overflow-hidden relative">
                    <div 
                      className={`h-full ${blockData.status === 'mev' ? 'bg-emerald-500/30' : 'bg-sky-500/30'}`}
                      style={{ 
                        width: `${Math.min(100, blockData.reward / 10000 * 100)}%`,
                        transition: 'width 0.5s ease-out' 
                      }}
                    ></div>
                    <div className="absolute inset-0 flex items-center justify-center text-xs text-indigo-200">
                      {blockData.reward > 0 ? `${blockData.reward.toLocaleString()} GWEI` : 'No Reward'}
                    </div>
                  </div>
                </div>
                
                {/* Add explanation */}
                <div className="mt-4 p-3 bg-indigo-900/30 rounded-lg text-xs text-indigo-200 leading-relaxed">
                  <p>Block rewards come from transaction fees paid by users. MEV-produced blocks may include extra rewards from Maximal Extractable Value strategies.</p>
                  {blockData.status === 'mev' && (
                    <p className="mt-1 text-emerald-200">This block was produced using MEV-Boost, which typically yields higher rewards.</p>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
